#!/usr/bin/env python3
"""Read/test only the disposable loopback gateway; no external model calls."""
import argparse
import json
import os
import time
import urllib.error
import urllib.request

ROOT = 'http://127.0.0.1:18081'
TEST_KEY = os.environ.get('EVERPLAIN_TEST_SERVICE_KEY', '')
HEADERS = {'Authorization':'Bearer ' + TEST_KEY,'Content-Type':'application/json','X-Session-Id':'synthetic-smoke-only'}

def call(path, body=None, auth=True):
    req = urllib.request.Request(ROOT+path,data=None if body is None else json.dumps(body).encode(),headers=HEADERS if auth else {},method='GET' if body is None else 'POST')
    try:
        with urllib.request.urlopen(req,timeout=15) as r:
            return r.status,dict(r.headers),r.read().decode()
    except urllib.error.HTTPError as e:
        return e.code,dict(e.headers),e.read().decode()

def main():
    if not TEST_KEY:
        raise SystemExit('Set EVERPLAIN_TEST_SERVICE_KEY to this disposable instance test key.')
    parser=argparse.ArgumentParser();parser.add_argument('--configured',action='store_true');args=parser.parse_args()
    assert call('/health',auth=False)[0]==200
    assert call('/v1/everplain/usage',auth=False)[0]==401
    for path in ['/api/v1/payment/config','/api/v1/admin/payment/config','/api/v1/admin/promo-codes','/api/v1/user/aff','/api/v1/subscriptions']:
        assert call(path,auth=False)[0]==404,path
    for path in ['/responses','/backend-api/codex/responses','/v1/images/generations']:
        assert call(path,{'model':'everplain-test'})[0]==404,path
    body={'model':'everplain-test','max_tokens':16,'messages':[{'role':'user','content':'synthetic benchmark'}]}
    status,headers,text=call('/v1/messages',body)
    if not args.configured:
        assert status==503 and json.loads(text)['error']['code']=='provider_not_configured',(status,text)
        print(json.dumps({'unconfigured':True,'sales_and_aliases_blocked':True,'anonymous_usage_rejected':True}));return
    assert status==200,(status,text)
    result=json.loads(text)
    assert result['usage']['input_tokens']==10 and result['usage']['output_tokens']==5,result
    assert headers.get('X-Everplain-Contract-Version')=='2026-10-04',headers
    body['stream']=True
    status,headers,text=call('/v1/messages',body)
    assert status==200 and 'message_stop' in text and 'synthetic ok' in text,(status,text)
    for _ in range(30):
        status,_,text=call('/v1/everplain/usage');data=json.loads(text)
        if status==200 and len(data['items'])>=2:break
        time.sleep(0.2)
    assert status==200 and len(data['items'])>=2,(status,text)
    assert all(item['model']=='everplain-test' and item['input_tokens']==10 for item in data['items']),data
    print(json.dumps({'unconfigured':False,'sync_generation':True,'stream_generation':True,'usage_export_count':len(data['items']),'model_mapping':data['items'][0].get('upstream_model'),'sales_and_aliases_blocked':True,'anonymous_usage_rejected':True}))

if __name__=='__main__':main()
