let ws;

function connect(){


    ws = new WebSocket("ws://192.168.0.2:8080/ws");

    ws.onopen = function (){
        
        console.log("connected to the Websocket server");
    }

    ws.onmessage = function(event) {

        console.log(JSON.parse(event.data))
        const {username, message} = JSON.parse(event.data) 
        let messageDisplay = document.getElementById("message");
        messageDisplay.innerHTML += `<p>USER:${username} says: ${message}</p>`;
    }

    ws.onclose = function () {
        console.log("WS connection closed... Retrying")
        setTimeout(connect, 1000)
    }

    ws.onerror = function (error){
        console.error(error)
    }

} 
function sendMessage(){
        let input = document.getElementById("messageInput");
        let message = input.value;
        input.value = "";
        let user = document.getElementById("userInput");
        let username = user.value;
        user.value = "";

        const payload = JSON.stringify({username: username, message: message})
        ws.send(payload)
    }
connect()