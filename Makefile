build-chatbot:
	./build/build.sh

deploy:
	./build/build.sh
	./deployments/cloudRun/deploy_chatbot.sh
