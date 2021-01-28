SHELL := /bin/bash

default: check-env prepare-secrets mount-s3 mount-s3-check publish-artifacts download-from-gh-check unmount-s3

check-env:
ifndef AWS_SECRET_ACCESS_KEY
	$(error AWS_SECRET_ACCESS_KEY is undefined)
endif
ifndef AWS_ACCESS_KEY
	$(error AWS_ACCESS_KEY is undefined)
endif
ifndef AWS_S3_BUCKET_NAME
	$(error AWS_S3_BUCKET_NAME is undefined)
endif

prepare-secrets:
	@echo "Generating secrets file into /etc/passwd-s3fs"
	@echo $(AWS_ACCESS_KEY):$(AWS_SECRET_ACCESS_KEY) > /etc/passwd-s3fs
	@chmod 600 /etc/passwd-s3fs

mount-s3:
	@echo "Mounting S3 into $(AWS_S3_MOUNT_DIRECTORY)"
	@s3fs $(AWS_S3_BUCKET_NAME) $(ARTIFACTS_DEST_FOLDER)

unmount-s3:
	@echo "Unmounting S3"
	@umount $(ARTIFACTS_DEST_FOLDER)

mount-s3-check:
	@echo "List files from s3 bucket to confirm mount"
	@ls -la $(AWS_S3_MOUNT_DIRECTORY)

download-from-gh-check:
	@echo "List all downloaded assets"
	@ls -la $(CURDIR)/assets

publish-artifacts:
	@echo "Publish artifacts"
	@/bin/publisher

.PHONY: prepare-secrets mount-s3 mount-s3-check publish-artifacts download-from-gh-check unmount-s3