all:
	@ cd server \
		&& govendor sync \
		&& cd ../ \
		&& docker build -t thequestion/test:$(IMAGE_TAG) .
