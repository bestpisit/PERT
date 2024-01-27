FROM nginx

# Copy Static file Folder into Nginx html hosted folder
COPY dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

#open the port
EXPOSE 80