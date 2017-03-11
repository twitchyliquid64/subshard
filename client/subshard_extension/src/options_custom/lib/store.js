//
// chrome.storage support by Jason Sterken
//
// store.js by Frank Kohlhepp
// Copyright (c) 2011 - 2012 Frank Kohlhepp
// https://github.com/frankkohlhepp/store-js
// License: MIT-license
//
(function () {
    var has = function (object, key) {
        return Object.prototype.hasOwnProperty.call(object, key);
    };

    var objectGetLength = function (object) {
        var count = 0;
        for (var key in object) {
            if (has(object, key)) { count++; }
        }

        return count;
    };

    var arrayIndexOf = function (array, item, from) {
        var length = array.length >>> 0;
        for (var i = (from < 0) ? Math.max(0, length + from) : from || 0; i < length; i++) {
            if (array[i] === item) { return i; }
        }

        return -1;
    };

    var arrayContains = function (array, item, from) {
        return arrayIndexOf(array, item, from) !== -1;
    };

    var arrayInclude = function (array, item) {
        if (!arrayContains(array, item)) { array.push(item); }
        return array;
    };

	var Store = function(storagetype, name, defaults, watcherSpeed, syncCallback) {
		this.name = name;
        this.defaults = defaults || {};
        this.watcherSpeed = watcherSpeed || 500;
        this.listeners = {};
		this.storage = storagetype;

		this.inSync = false;

		var self = this;
		this.storage.get(null, function(items) {
			var name = "store." + self.name + ".";
			for(var key in items) {
				if(key.substring(0, name.length) === name) {
					var thisKey = key.substring(name.length);
					var value = items[key];
					if(value == self.localFlag(thisKey)) {
						try {
							value = JSON.parse(localStorage.getItem(self.localKey(thisKey)));
						} catch(e) {
							value = null;
						}
						if(value === null && self.defaults.hasOwnProperty(thisKey)) {
							//set default so full pass doesn't overwrite flag
							value = self.defaults[thisKey];
						}
						self.set(thisKey, value, true, true); //name, value, localOnly, don't queue

					} else {
						self.set(thisKey, value, false, true);
					}
				}
			}
			self.inSync = true;

			//applyDefaults synchronously
			for(var defaultKey in self.defaults) {
				var check = self.stateGet(defaultKey);
				if(check === null || typeof check == "undefined") {
					self.set(defaultKey, self.defaults[defaultKey]);
				}
			}

			if(typeof syncCallback == "function") {
				syncCallback(this);
			}
		});
	};

	this.Store = (function() {
		return (function(name, defaults, watcherSpeed, syncCallback) {
			return new Store(chrome.storage.local, name, defaults, watcherSpeed, syncCallback);
		})
	})();

	this.StoreSync = (function() {
		return (function(name, defaults, watcherSpeed, syncCallback) {
			return new Store(chrome.storage.sync, name, defaults, watcherSpeed, syncCallback);
		})
	})();

    Store.clear = function () {
		this.storage.clear();
		localStorage.clear();
    };

	Store.prototype.stateGet = function(key) {
		try {
			return JSON.parse(localStorage.getItem("store." + this.name + "_state." + key));
		}catch(e){
			return null;
		}
	};

	Store.prototype.stateSet = function(key, value) {
		try {
			localStorage.setItem("store." + this.name + "_state." + key, JSON.stringify(value));
		}catch(e){}
	};

	Store.prototype.stateReset = function() {
		var prefix = "store." + this.name + "_state.";
		var keys = Object.keys(localStorage);
		for(var i=0; i < keys.length; i++) {
			var key = keys[i];
			if(key.substring(0, prefix.length) === prefix) {
				localStorage.removeItem(key);
			}
		}
	};
	Store.prototype.stateRemove = function(key) {
		localStorage.removeItem("store." + this.name + "_state." + key);
	};

	Store.prototype.queueReady = true;
	Store.prototype.queue = function(key, value) {
		var qName = "store."+ this.name +"_queue";
		var q = localStorage.getItem(qName);
		if(q === null) {
			q = {}
		} else {
			try {
				q = JSON.parse(q);
			}catch(e) { q = {}; }
		}

		q[key] = value;
		localStorage.setItem(qName, JSON.stringify(q));

		if(this.queueReady === true) {
			this.queueWrite();
		}
	}

	Store.prototype.queueCommit = function() {
		this.queueReady = true;
		this.queueWrite();
	}
	Store.prototype.queueWrite = function() {
		var qName = "store."+ this.name +"_queue";
		var q = localStorage.getItem(qName);
		if(q === null) {
			q = {}
		} else {
			try {
				q = JSON.parse(q);
			}catch(e) { q = {}; }
		}

		if(Object.keys(q).length > 0) {
			this.storage.set(q);
			localStorage.removeItem(qName);

			this.queueReady = false;
			var self = this;
			this.queueReady = setTimeout(function() {
				self.queueReady = true;
				self.queueWrite();
			}, 5000);
		}
	}

    Store.prototype.applyDefaults = function () {
		var self = this;
		var prefix = "store." + this.name + ".";

		var getKeys = [];
		var baseKeys = [];
        for (var key in this.defaults) {
			baseKeys.push(key);
			getKeys.push(prefix + key);
		}
		this.storage.get(getKeys, function(items) {
			for(var i=0; i < getKeys.length; i++) {
				if(!items.hasOwnProperty(getKeys[i])) {
					self.set(baseKeys[i], self.defaults[baseKeys[i]]);
				}
			}
		});

        return this;
    };

    Store.prototype.watcher = function (force) {
        if (this.watcherTimer) {
            clearTimeout(this.watcherTimer);
        }

        if (objectGetLength(this.listeners) || force) {
            this.newObject = this.toObject();

            if (this.oldObject) {
                for (var key in this.newObject) {
                    if (has(this.newObject, key) && this.newObject[key] !== this.oldObject[key]) {
                        this.fireEvent(key, this.newObject[key]);
                    }
                }

                for (var key in this.oldObject) {
                    if (has(this.oldObject, key) && !has(this.newObject, key)) {
                        this.fireEvent(key, this.newObject[key]);
                    }
                }
            }

            this.oldObject = this.newObject;
            var that = this;
            this.watcherTimer = setTimeout(function () {
                that.watcher();
            }, this.watcherSpeed);
        }

        return this;
    };

    Store.prototype.get = function (key, callback) {
		if(this.inSync) {
			var value = this.stateGet(key);

			if(value == this.localFlag(key)) {
				try {
					value = JSON.parse(localStorage.getItem(this.localKey(key)));
				} catch(e) {
					value = null;
				}
			}

			if(typeof callback == "undefined") {
				return value;
			} else {
				callback(value);
			}
		} else {
			var self = this;
			var fullKey = "store." + this.name + "." + key;
			this.storage.get(fullKey, function(items) {
				var value = items[fullKey];
				if(value == self.localFlag(key)) {
					try {
						value = JSON.parse(localStorage.getItem(self.localKey(key)));
					} catch(e) {
						value = null;
					}
				}

				self.stateSet(key, value);

				if(typeof callback == "function") {
					callback(value);
				}
			});
		}
    };

    Store.prototype.set = function (key, value, localOnly, init) {
        if (value === undefined) {
            this.remove(key);
        } else {
            if (typeof value === "function") { value = null; }

			this.stateSet(key, value);

			if(localOnly) {
				var jsonValue = jsonValue = JSON.stringify(value);
				localStorage.setItem(this.localKey(key), jsonValue);
				value = this.localFlag(key);
			}

			if(!init) {
				this.queue("store." + this.name + "." + key, value);
			}
        }

        return this;
    };

    Store.prototype.remove = function (name) {
        this.storage.remove("store." + this.name + "." + name);
		this.stateRemove(name);
        return this.applyDefaults();
    };

    Store.prototype.reset = function () {
		var self = this;
		this.storage.get(null, function(items) {
			var name = "store." + self.name + ".";
			var remove = [];
			for(var key in items) {
				if(key.substring(0, name.length) === name) {
					remove.push(key);
				}
			}
			self.storage.remove(remove);
			self.applyDefaults();
		});
		this.stateReset();

		return this;
    };

    Store.prototype.toObject = function (callback) {
		if(this.inSync) {
			var keys = Object.keys(localStorage);
			var retObj = {};
			var prefix = "store." + this.name + "_state.";
			for(var i=0; i < keys.length; i++) {
				var key = keys[i];
				if(key.substring(0, prefix.length) === prefix) {
					retObj[key.substring(prefix.length)] = localStorage.getItem(key);
				}
			}

			if(typeof callback == "undefined") {
				return retObj;
			} else {
				callback(retObj);
			}
		} else {
			var self = this;
			this.storage.get(null, function(items) {
				var values = {};
				var name = "store." + self.name + ".";
				for(var key in items) {
					if(key.substring(0, name.length) === name) {
						var thisKey = key.substring(name.length);
						var value = items[key];
						if (value !== undefined) { values[thisKey] = value; }
					}
				}

				callback(values);
			});
		}
    };

    Store.prototype.fromObject = function (values, merge) {
        if (!merge) { this.reset(); }
        for (var key in values) {
            if (has(values, key)) {
                this.set(key, values[key]);
            }
        }

        return this;
    };

	Store.prototype.localFlag = function(key) {
		return "store." + this.name + "." + key + "___CHECK_LOCAL___";
	};
	Store.prototype.localKey = function(key) {
		return "store." + this.name + "_localonly." + key;
	};

    Store.prototype.addEvent = function (selector, callback) {
        this.watcher(true);
        if (!this.listeners[selector]) { this.listeners[selector] = []; }
        arrayInclude(this.listeners[selector], callback);
        return this;
    };

    Store.prototype.removeEvent = function (selector, callback) {
        for (var i = (this.listeners[selector].length - 1); i >= 0; i--) {
            if (this.listeners[selector][i] === callback) { this.listeners[selector].splice(i, 1); }
        }

        if (!this.listeners[selector].length) { delete this.listeners[selector]; }
        return this;
    };

    Store.prototype.fireEvent = function (name, value) {
        var selectors = [name, "*"];
        for (var i = 0; i < selectors.length; i++) {
            var selector = selectors[i];
            if (this.listeners[selector]) {
                for (var j = 0; j < this.listeners[selector].length; j++) {
                    this.listeners[selector][j](value, name, this.name);
                }
            }
        }

        return this;
    };
}());
