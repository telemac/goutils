# com.plugis.browser cloud events

Topic : com.plugis.browser

Cloud events : 
- Type : com.plugis.browser.open
- parameters : {"url": "https://www.google.fr"}
- sample : nats-pub -s 'https://nats1.plugis.com' 'com.plugis.browser' '{"type": "com.plugis.browser.open","data": {"url": "https://www.google.fr"}, "id": "123","source": "manual","specversion": "1.0"}'

