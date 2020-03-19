.PHONY: build
build:
	yarn build

deploy: build
	#aws s3 cp --profile delphis --region us-west-2 ./build s3://com.delphishq.staging.app/ --recursive --exclude "fonts/*"
	aws s3 sync --profile delphis --region us-west-2 --acl public-read ./build s3://app-staging.delphishq.com/ --exclude "fonts/*" --delete
	aws s3 cp --profile delphis --region us-west-2 --acl public-read ./build/index.html s3://app-staging.delphishq.com/index.html --cache-control max-age=0,no-cache,no-store,must-revalidate