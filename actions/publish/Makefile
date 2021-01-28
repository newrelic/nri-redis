SHELL := /bin/bash

AWS_S3_MOUNT_DIRECTORY ?= /mnt/s3

default: check-env prepare-secrets mount-s3 mount-s3-check upload-artifacts download-from-gh-check

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
	@s3fs $(AWS_S3_BUCKET_NAME) $(AWS_S3_MOUNT_DIRECTORY)

mount-s3-check:
	@echo "List files from s3 bucket to confirm mount"
	@ls -la $(AWS_S3_MOUNT_DIRECTORY)

download-from-gh-check:
	@echo "List all downloaded assets"
	@ls -la $(CURDIR)/assets

upload-artifacts:
	@echo "Upload artifacts"
	@./uploader

# jmx
# repos
# nrjmx_1.5.2-1_amd64.deb => infrastructure_agent/linux/apt/pool/main/n/nrjmx/
# nrjmx-1.5.2-1.x86_64.rpm => infrastructure_agent/linux/yum/el/5/x86_64/
# nrjmx-1.5.2-1.x86_64.rpm => infrastructure_agent/linux/yum/el/6/x86_64/
# nrjmx-1.5.2-1.x86_64.rpm => infrastructure_agent/linux/yum/el/7/x86_64/
# nrjmx-1.5.2-1.x86_64.rpm => infrastructure_agent/linux/yum/el/8/x86_64/
# nrjmx-1.6.2-1.x86_64.rpm => infrastructure_agent/linux/zypp/sles/11.4/x86_64/
# nrjmx-1.6.2-1.x86_64.rpm => infrastructure_agent/linux/zypp/sles/12.1/x86_64/
# nrjmx-1.6.2-1.x86_64.rpm => infrastructure_agent/linux/zypp/sles/12.2/x86_64/
# nrjmx-1.6.2-1.x86_64.rpm => infrastructure_agent/linux/zypp/sles/12.3/x86_64/
# nrjmx-1.6.2-1.x86_64.rpm => infrastructure_agent/linux/zypp/sles/12.4/x86_64/

.PHONY: prepare-secrets mount-s3 mount-s3-check upload-artifacts download-from-gh-check