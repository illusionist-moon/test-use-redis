var https = require('https');
var fs = require('fs');

var options = {
    key: fs.readFileSync('./server-key.pem'),
    ca: [fs.readFileSync('./ca-cert.pem')],
    cert: fs.readFileSync('./server-cert.pem'),
    // pfx: fs.readFileSync('./server.pfx'),
    // passphrase: '20020911',
}

https.createServer(options,function(req,res){
	res.writeHead(200);
	res.end('hello world\n');
}).listen(8080,'127.0.0.1');