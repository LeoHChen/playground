var W3CWebSocket = require('websocket').w3cwebsocket;
let client = new W3CWebSocket("wss://ws.s0.b.hmny.io");
// let client = new W3CWebSocket("wss://ws.s0.t.hmny.io");
client.onopen = function() {
       console.log('WebSocket Client Connected');
};
client.onerror = function() {
       console.log('Connection Error');
};
