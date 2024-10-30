echo "Note you must be running this on wsl if you are on a windows machine"
echo "To install wsl type in wsl --install"

sudo apt install make
Y

sudo apt install build-essential
Y

sudo apt install tar
Y

sudo apt install gzip
Y

sudo apt install flex
Y

sudo apt install m4
Y
PATH=$PATH:/usr/local/m4/bin/

wget http://ftp.gnu.org/gnu/bison/bison-2.3.tar.gz
tar -xvzf bison-2.3.tar.gz
./configure --prefix=/usr/local/bison --with-libiconv-prefix=/usr/local/libiconv/
make
sudo make install
cd ..

sudo apt install perl
Y

sudo apt install zlib1g
Y

sudo apt-get install libicu-dev
Y

sudo apt-get install zlib1g-dev
Y

sudo apt install pkgconf
Y

git clone https://git.postgresql.org/git/postgresql.git
cd postgresql/
./configure --without-readline
make
su
make install
adduser postgres