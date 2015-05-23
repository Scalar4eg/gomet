var socket = io()

$(function() {

    var recv_types = [
        "auth_response",
         "new_contact",
         "delete_contact",
         "message_recv",
         "message_accepted",
         "message_read",
         "contact_status",
    ]

    for (var i in recv_types) {
        var t = recv_types[i];

        (function(type) {
            socket.on(type, function(data) {
                console.log("Got event: "+type, JSON.parse(data))
            })
        })(t);
    }

    window.send = function(message_type, message_data) {
        socket.emit(message_type, JSON.stringify(message_data))
    }
});