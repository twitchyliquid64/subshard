function corsGlobal() {
  console.log('corsGlobal');
  chrome.runtime.sendMessage({
    'msg': 'toggleCORSGlobal',
  }, function(){
    updateUI();
  });
}

function corsSensitive() {
  console.log('corsSensitive');
  chrome.runtime.sendMessage({
    'msg': 'toggleCORSSensitive',
  }, function(){
    updateUI();
  });
}

function userAgentHide(){
  console.log('toggleMockUserAgent');
  chrome.runtime.sendMessage({
    'msg': 'toggleMockUserAgent',
  }, function(){
    updateUI();
  });
}


function updateUI(){
  var cross_resources_disabled = chrome.extension.getBackgroundPage().cross_resources_disabled;
  var sensitive_protection_disabled = chrome.extension.getBackgroundPage().sensitive_protection_disabled;
  var hide_user_agent = chrome.extension.getBackgroundPage().hide_user_agent;

  if (cross_resources_disabled) {
    document.getElementById('corsGlobal').innerHTML = 'Enable CORS globally';
  } else {
    document.getElementById('corsGlobal').innerHTML = 'Disable CORS globally';
  }

  if (sensitive_protection_disabled) {
    document.getElementById('corsSensitive').innerHTML = 'Disable CORS on sensitive sites (recommended)';
    document.getElementById('warn').innerHTML = 'WARNING: Sensitive domains are allowed to make cross-domain requests. You should disable this.'
  } else {
    document.getElementById('corsSensitive').innerHTML = 'Enable CORS on sensitive sites (!!!)';
    document.getElementById('warn').innerHTML = '';
  }

  if (hide_user_agent){
    document.getElementById('hide_user_agent').innerHTML = 'Set User-Agent to correct value';
  } else {
    document.getElementById('hide_user_agent').innerHTML = 'Set User-Agent to IE11';
  }
}


document.addEventListener('DOMContentLoaded', function() {
  document.getElementById('corsGlobal').addEventListener('click', corsGlobal);
  document.getElementById('corsSensitive').addEventListener('click', corsSensitive);
  document.getElementById('hide_user_agent').addEventListener('click', userAgentHide);
  updateUI();
});
