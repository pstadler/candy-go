# Candy on Go

Deploy [Candy](http://candy-chat.github.com/candy/) as a [Go](http://golang.org) application.

## Setup

* Check out the submodules: `git submodule update --init`
* Copy `config.json.sample` to `config.json`
* Change the preferences in `config.json` to fit your setup (http://candy-chat.github.io/candy/#configuration)
* Run it with `make`

## Heroku

* Copy `.env.sample` to `.env`
* Change the preferences in `.env` to fit your setup (http://candy-chat.github.io/candy/#configuration)
* `heroku create -b https://github.com/kr/heroku-buildpack-go.git`
* `heroku plugins:install git://github.com/ddollar/heroku-config.git`
* `heroku config:push` to push your `.env` configuration
* `git push heroku master` to deploy it
* `heroku apps:open` to enjoy it
