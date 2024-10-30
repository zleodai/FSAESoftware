echo "Note you must be running this on wsl if you are on a windows machine"
echo "To install wsl run: wsl --install -d Ubuntu"

sudo apt-get install make -y
sudo apt-get install build-essential -y
sudo apt-get install tar -y
sudo apt-get install gzip -y
sudo apt-get install flex -y
sudo apt-get install m4 -y
PATH=$PATH:/usr/local/m4/bin/

wget http://ftp.gnu.org/gnu/bison/bison-2.3.tar.gz
tar -xvzf bison-2.3.tar.gz
cd bison-2.3
./configure --prefix=/usr/local/bison --with-libiconv-prefix=/usr/local/libiconv/
make
sudo make install
cd ..
rm -r bison-2.3
rm bison-2.3.tar.gz

sudo apt-get install bison -y
sudo apt-get install perl -y
sudo apt-get install zlib1g -y
sudo apt-get install libicu-dev -y
sudo apt-get install zlib1g-dev -y
sudo apt-get install pkgconf -y

git clone https://git.postgresql.org/git/postgresql.git
cd postgresql/
./configure --without-readline
make
su
make install
adduser postgres

cd ..