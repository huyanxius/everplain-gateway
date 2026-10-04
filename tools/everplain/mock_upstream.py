#!/usr/bin/env python3
"""Loopback-only Anthropic fixture; never forwards or reads real credentials."""
import argparse
import ssl
import json
import time
import uuid
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

class Handler(BaseHTTPRequestHandler):
    def log_message(self, *_):
        pass

    def do_POST(self):
        length = int(self.headers.get('Content-Length', '0'))
        if length > 4 * 1024 * 1024:
            self.send_error(413)
            return
        body = json.loads(self.rfile.read(length))
        if self.path.split('?')[0] != '/v1/messages':
            self.send_error(404)
            return
        # Constant artificial network latency makes small-concurrency comparisons
        # reproducible. This is not a claim about any real provider's latency.
        time.sleep(0.025)
        rid = 'synthetic-' + uuid.uuid4().hex
        model = body.get('model', 'claude-sonnet-4-20250514')
        self.send_response(200)
        self.send_header('request-id', rid)
        if body.get('stream'):
            self.send_header('Content-Type', 'text/event-stream')
            self.end_headers()
            events = [
                ('message_start', {'type': 'message_start', 'message': {'id':rid,'type':'message','role':'assistant','model':model,'content':[],'usage':{'input_tokens':10,'output_tokens':0}}}),
                ('content_block_start', {'type':'content_block_start','index':0,'content_block':{'type':'text','text':''}}),
                ('content_block_delta', {'type':'content_block_delta','index':0,'delta':{'type':'text_delta','text':'synthetic ok'}}),
                ('content_block_stop', {'type':'content_block_stop','index':0}),
                ('message_delta', {'type':'message_delta','delta':{'stop_reason':'end_turn'},'usage':{'output_tokens':5}}),
                ('message_stop', {'type':'message_stop'}),
            ]
            for event, data in events:
                self.wfile.write(('event: ' + event + '\ndata: ' + json.dumps(data) + '\n\n').encode())
                self.wfile.flush()
        else:
            output = json.dumps({'id':rid,'type':'message','role':'assistant','model':model,'content':[{'type':'text','text':'synthetic ok'}],'stop_reason':'end_turn','usage':{'input_tokens':10,'output_tokens':5,'cache_creation_input_tokens':0,'cache_read_input_tokens':0}}).encode()
            self.send_header('Content-Type','application/json')
            self.send_header('Content-Length',str(len(output)))
            self.end_headers()
            self.wfile.write(output)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--tls-cert')
    parser.add_argument('--tls-key')
    args = parser.parse_args()
    if bool(args.tls_cert) != bool(args.tls_key):
        parser.error('--tls-cert and --tls-key must be supplied together')
    server = ThreadingHTTPServer(('127.0.0.1',18082),Handler)
    if args.tls_cert:
        context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        context.load_cert_chain(args.tls_cert, args.tls_key)
        server.socket = context.wrap_socket(server.socket, server_side=True)
    server.serve_forever()
