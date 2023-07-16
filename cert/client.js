var https = require('https')
var fs = require('fs')

var options = {
    hostname: '127.0.0.1',
    port: 8080,
    path: '/',
    method: 'GET',
    // pfx: fs.readFileSync('./server.pfx'),
    // passphrase: '20020911',
    key: fs.readFileSync('./client-key.pem'),
    ca: [fs.readFileSync('./ca-cert.pem')],
    cert: fs.readFileSync('./client-cert.pem'),
    agent: false,
}

options.agent = new https.Agent(options)
var req = https.request(options, function (res) {
    console.log('statusCode: ', res.statusCode)
    console.log('headers: ', res.headers)
    res.setEncoding('utf-8')
    res.on('data', function (d) {
        console.log(d)
    })
})

req.end()

req.on('error', function (e) {
    console.log(e)
})
