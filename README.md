# httpdigest
A simple Go(lang) drop-in replacement for "net/http/Client" with Digest Auth capability.

## Warning
The API design is on going. Watch out when updating.

## Usage
    import (
        "net/http"
        
        "github.com/iwat/httpdigest"
    )
    
    func test() {    
        type IdentStore struct {
            user string
            pass string
        }
    
        client := httpdigest.Client{AuthHandler: identStore}
    
        request, err := http.NewRequest("GET", url, nil)
    
        if err != nil {
            panic(err)
        }
        
        resp, err := client.Do(request)
    
        if err != nil {
            panic(err)
        }
    }

    func (id IdentStore) HandleAuth(resp *http.Response, req *http.Request) {
        challenge := httpdigest.ChallengeFromResponse(resp)
        challenge.ApplyAuth(id.user, id.pass, req)
    }

## Legal
This application is developed under MIT license, and can be used for open and proprietary projects.

Copyright (c) 2015 Chaiwat Shuetrakoonpaiboon

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
