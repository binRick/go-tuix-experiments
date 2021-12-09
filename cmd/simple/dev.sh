nodemon -w . -I -e go --delay .2 -x sh -- -c 'reset;reap passh go run .||true'
