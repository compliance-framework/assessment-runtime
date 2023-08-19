# Use an official Node.js runtime as the parent image
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install the application's dependencies inside the container
RUN npm install

# Bundle the application source inside the container
COPY . .

# Specify the port the app runs on
EXPOSE 3000

# Define the command to run the app
CMD ["node", "server.js"]
