let pollXhr = new XMLHttpRequest();
let pollSince = Date.now() - 86400000
const longPollTimeout = 30

pollNow = () => {
    pollXhr.open('GET', `/events?timeout=${longPollTimeout}&category=refresh&since_time=${pollSince}`);
    pollXhr.send();
}

pollXhr.onload = () => {
    if (pollXhr.status >= 200 & pollXhr.status < 300) {
        let pollData = JSON.parse(pollXhr.responseText);
        if (typeof pollData.events != 'undefined') {
            pollSince = pollData.events[0].timestamp;
            xhr.open('POST', '/graphql');
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.send(JSON.stringify(body));
            }
        if (pollData.timeout != undefined) {
            pollSince = pollData.timestamp;
        }
        pollNow();
    } else {
        // TODO: Better error handling. For now just passing the xhr object to a function which usually gets a ProgressEvent
        pollError(pollXhr);
    }
}

pollXhr.onerror = e => pollError(e);

pollError = (e) => {
    console.error("Long poll returned error");
    console.log(e);
    const errorDiv = document.getElementById('error');
    const errMsg = document.createElement('p');
    errMsg.innerText = `Polling error. Live updates have stopped. Correct and refresh page.`
    errorDiv.appendChild(errMsg);

}

pollNow();