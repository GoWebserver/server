openssl req -newkey rsa:8192 \
  -new -nodes -x509 \
  -days 7300 \
  -out cert.crt \
  -keyout key.key \
  -subj "/C=DE/ST=Bayern/L=Donauwörth/O=H3rmts in Dark/OU=Unit/CN=localhost"