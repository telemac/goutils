class CloudEvents {
    constructor(params) {
        console.log("CloudEvents.constructor", params)
    }

    send (type, data, topic) {
        console.log("CloudEvents.send", type, data, topic)
        let event = {
            topic,
            type,
            request: true,
            timeout: 5,
            source: 'web:'+window.location.href,
            specversion: '1.0',
            datacontenttype: 'application/json',
            data: data
        }

        return fetch('/cloudevents/send', {
            method: 'POST',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify(event)
        })
            .then(response => {
                console.log('response', response);response
            } )
            .catch(err => console.error(err))
    }
}
