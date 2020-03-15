worker_processes  1;

events {
    worker_connections 1024;
}

http {
    server {
        listen 8000;
        server_name local.delphishq.com;

        location = /graphiql {
	      #proxy_pass http://local.delphishq.com:8080;
	      proxy_pass http://staging.delphishq.com;
        }

        location = /query {
	      #proxy_pass http://local.delphishq.com:8080;
	      proxy_pass http://staging.delphishq.com;
        }

        location /twitter {
	      #proxy_pass http://local.delphishq.com:8080;
	      proxy_pass http://staging.delphishq.com;
        }

        location /sockjs-node {
            proxy_set_header X-Real-IP  $remote_addr;
            proxy_set_header X-Forwarded-For $remote_addr;
            proxy_set_header Host $host;

            proxy_pass http://local.delphishq.com:3000; 

            proxy_redirect off;

            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }

        location / {
            proxy_pass http://local.delphishq.com:3000;
        }

    }
}
