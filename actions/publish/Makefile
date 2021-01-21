all: prepare-secrets mount-s3 mount-s3-check download-from-gh download-from-gh-check

prepare-secrets:
	@echo "Generating secrets file into /etc/passwd-s3fs"
	@echo $(AWS_ACCESS_KEY):$(AWS_SECRET_ACCESS_KEY) > /etc/passwd-s3fs
	@chmod 600 /etc/passwd-s3fs

mount-s3:
	@echo "Mounting S3 into $(AWS_S3_MOUNT_DIRECTORY)"
	@s3fs $(AWS_S3_BUCKET_NAME) $(AWS_S3_MOUNT_DIRECTORY)

mount-s3-check:
	@echo "List files from s3 bucket to confirm mount"
	@ls -la $(AWS_S3_MOUNT_DIRECTORY)

download-from-gh:
	@echo "Download packages from github release assets"
	@bash $(CURDIR)/scripts/download-gh-assets.sh

download-from-gh-check:
	@echo "List all downloaded assets"
	@ls -la $(CURDIR)/assets

.PHONY: prepare-secrets mount-s3 mount-s3-check download-from-gh download-from-gh-check