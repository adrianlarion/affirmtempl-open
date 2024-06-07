
# about this project
I am passionate about self development and programming. Thus I created an affirmations app with Go - which is becoming quickly my favourite language.
The code is by no means perfect, although I tried to keep it as clean as possible.
I've used Echo and Templ on backend, materialize.css and jQuery on frontend.

Feel free to modify it and improve it as you see fit.


![](./docs/screenshot1.png)


# tooling needed to be installed before running
```
apt-get install xdotool
apt install mariadb-server


apt-get install npm
npm install -g sass
npm install -g browserify

cd internal/frontend
npm install

go mod tidy
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest

```

- also install the migrate tool.
- https://github.com/golang-migrate/migrate/releases
- You can download the .deb package from assets and double click to install


# mysql setup:
```
sudo mysql
```
- then:
```
CREATE DATABASE affirmtempl;
CREATE USER 'affirmtempl'@'localhost';
GRANT ALL PRIVILEGES ON affirmtempl.* TO 'affirmtempl'@'localhost' WITH GRANT OPTION;
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'affirmtempl'@'localhost' IDENTIFIED BY 'pass';

```

- then migrate
```
migrate -path=./migrations -database="mysql://affirmtempl:pass@tcp(localhost)/affirmtempl" up

```


# run dev
1. run `air` in the root directory of this project to run a development server and live reload.
2. Open 'localhost:4000'  in your browser to use this app.


- the reload script will try to find a Google Chrome window for hot reloading
- in the script located at `./scripts/devreload.sh`
`WID=$(xdotool search --name "Google Chrome")
`
- If you have Firefox or another browser, change the name to the respective browser instead of "Google Chrome"


 
# style modifications
* static/sass/app/_appstyle.scss
* static/sass/materialize.scss
* static/sass/_color-variables.scss
* static/sass/components/_variables.scss
* static/js/app.js



# licensing and disclaimer
I'm releasing this under a Creative Common CC0 license. Basically you can do whatever you want with it, no need for attribution either.
https://creativecommons.org/public-domain/cc0/

Disclaimer:
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
