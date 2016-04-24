First install the NGINX block in this directory. (There is an experimental Makefile that will do this, just try `sudo make install`.

To use letsencrypt follow these steps:

```
git clone https://github.com/letsencrypt/letsencrypt
cd letsencrypt
sudo service nginx stop
sudo ./letsencrypt-auto certonly --standalone --email youremail@somewhere.com -d yourserver.com
sudo service nginx start
```

Then startup `cowyo` with

```bash
sudo ./cowyo -p :8001 -key /etc/letsencrypt/live/yourserver.com/privkey.pem -crt /etc/letsencrypt/live/yourserver.com/cert.pem yourserver.com
```
