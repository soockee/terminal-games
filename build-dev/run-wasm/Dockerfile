FROM nginx:stable-alpine

WORKDIR /usr/share/nginx/html

COPY index.html index.js ldtk-snake.wasm wasm_exec.js /usr/share/nginx/html/

RUN rm -rf nginx.conf

COPY nginx.conf .

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]