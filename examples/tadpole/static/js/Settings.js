var Settings = function () {
    this.protocol = "ws:";
    this.host = window.location.hostname;
    this.port = 8081;
    this.path = '/ws';
    if (window.location.protocol === 'https:') {
        this.protocol = "wss:";
        if (window.location.port === "") {
            this.port = 443;
        }
    }
};
