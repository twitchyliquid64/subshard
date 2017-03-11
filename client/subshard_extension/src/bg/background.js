// if you checked "fancy-settings" in extensionizr.com, uncomment this lines

// var settings = new Store("settings", {
//     "sample_setting": "This is how you use Store.js to remember values"
// });

var allowedCrossDomainPages = {};
var allowedCrossDomainTargets = {
  'developer.cdn.mozilla.net': {},
};
var protectedHostSuffixes = [
  '.onion',
];

var tabInfo = {};
var cross_resources_disabled = false;
var sensitive_protection_disabled = false;
var hide_user_agent = false;

chrome.tabs.onUpdated.addListener(function(tabId, changeInfo, tab) {
  tabInfo[tabId] = tab;
});
chrome.tabs.onActivated.addListener(function(tabId, changeInfo, tab) {
  tabInfo[tabId] = tab;
});
chrome.tabs.onRemoved.addListener(function(tabId, removeInfo) {
  delete tabInfo[tabId];
});







function isProtectedOnionLink(url){
  if (sensitive_protection_disabled){
    return false;
  }
  for (var i = 0; i < protectedHostSuffixes.length; i++){
    if (url.host.endsWith(protectedHostSuffixes[i])){
      return true;
    }
  }
  return false;
}

function checkURLMatch(entry, url){
  if (!entry){
    return false;
  }
  if (entry.urls){
    for (var i = 0; i < entry.urls.length; i++){
      if (entry.urls[i] == url.pathname)return true;
    }
  } else {
    return true; //There is an entry for this domain but no URLS, whitelist whole domain
  }
  return false;
}

//check our pages for cross-domain requests, abort them if they do not match the whitelists.
chrome.webRequest.onBeforeRequest.addListener(
  function(details) {
    if (tabInfo[details.tabId]) {
      var tab_url = new URL(tabInfo[details.tabId].url);
      var request_url = new URL(details.url);
      var should_protect = !!(isProtectedOnionLink(tab_url) || cross_resources_disabled);

      if (should_protect && (tab_url.host != request_url.host)){
        console.log("ABORT REQUEST CANDIDATE: ", tab_url.host, '-->', request_url.host, cross_resources_disabled);
        if (checkURLMatch(allowedCrossDomainPages[tab_url.host], tab_url.pathname))
          return {};
        if (checkURLMatch(allowedCrossDomainTargets[request_url.host], request_url.pathname))
          return {};
        console.log("Aborting request: ", request_url.href);
        return {cancel: true};
      }
    }
  },
  {urls: ["<all_urls>"]},
  ["blocking"]
);



// Strip the referer header
chrome.webRequest.onBeforeSendHeaders.addListener(onHeadersIntercept, {urls: ["<all_urls>"]}, ["blocking", "requestHeaders"]);
function onHeadersIntercept(details){
  for (var i = 0; i < details.requestHeaders.length; ++i) {
    if (details.requestHeaders[i].name === 'Referer') {
      details.requestHeaders.splice(i, 1);
      break;
    }
    if (hide_user_agent && details.requestHeaders[i].name === 'User-Agent'){
      details.requestHeaders[i].value = "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko";
    }
  }
  return {requestHeaders: details.requestHeaders};
}




chrome.storage.onChanged.addListener(function(changes, namespace) {
  if (namespace == 'sync'){
    for (key in changes) {
      console.log("Data change to: " + key);
      var storageChange = changes[key];
      switch (key){
        case 'cross_resources_disabled':
          cross_resources_disabled = storageChange.newValue;
          break;
        case 'sensitive_protection_disabled':
          sensitive_protection_disabled = storageChange.newValue;
          break;
        case 'allowedCrossDomainPages':
          allowedCrossDomainPages = storageChange.newValue;
          break;
        case 'allowedCrossDomainTargets':
          allowedCrossDomainTargets = storageChange.newValue;
          break;
        case 'hide_user_agent':
          hide_user_agent = storageChange.newValue;
      }
    }
  }
});

function loadStorage(){
  chrome.storage.sync.get(['cross_resources_disabled', 'sensitive_protection_disabled', 'allowedCrossDomainTargets', 'allowedCrossDomainPages', 'hide_user_agent'], function(data){
    cross_resources_disabled = data['cross_resources_disabled'] || false;
    sensitive_protection_disabled = data['sensitive_protection_disabled'] || false;
    allowedCrossDomainPages = data['allowedCrossDomainPages'] || allowedCrossDomainPages;
    allowedCrossDomainTargets = data['allowedCrossDomainTargets'] || allowedCrossDomainTargets;
    protectedHostSuffixes = data['protectedHostSuffixes'] || protectedHostSuffixes;
    hide_user_agent = data['hide_user_agent'] || false;
    console.log("Storage loaded:", data);
  })
}

loadStorage();


function save(){
  var settings = {
    'cross_resources_disabled': cross_resources_disabled,
    'sensitive_protection_disabled': sensitive_protection_disabled,
    'allowedCrossDomainPages': allowedCrossDomainPages,
    'allowedCrossDomainTargets': allowedCrossDomainTargets,
    'protectedHostSuffixes': protectedHostSuffixes,
    'hide_user_agent': hide_user_agent,
  };
  chrome.storage.sync.set(settings, function() {
    console.log("Saved.")
  });
}




//Message handler from the UI
chrome.runtime.onMessage.addListener(function(message)  {
    console.log("Got message:", message);
    if (message.msg == "toggleCORSGlobal"){
      cross_resources_disabled = !cross_resources_disabled;
      save();
    }
    if (message.msg == "toggleCORSSensitive"){
      sensitive_protection_disabled = !sensitive_protection_disabled;
      save();
    }
    if (message.msg == "toggleMockUserAgent"){
      hide_user_agent = !hide_user_agent;
      save();
    }
});
