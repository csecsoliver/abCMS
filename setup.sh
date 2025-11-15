sudo apt-get -y install --no-install-recommends wget gnupg ca-certificates lsb-release
wget -O - https://openresty.org/package/pubkey.gpg | sudo gpg --dearmor -o /usr/share/keyrings/openresty.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/openresty.gpg] http://openresty.org/package/ubuntu $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/openresty.list > /dev/null
sudo apt-get update
sudo apt-get -y install openresty
sudo apt install -y uuid-dev
sudo apt-get install libmagickwand-dev -y
sudo apt install luarocks -y
sudo luarocks install lapis
sudo luarocks install moonscript
sudo luarocks install markdown
sudo luarocks install lsqlite3
sudo luarocks install bcrypt
sudo luarocks install lua-uuid
sudo luarocks install luafilesystem
sudo luarocks install magick
