SHELL := /bin/bash
DOWNLOAD_SCRIPT ?= $(CURDIR)/scripts/download-gh-assets.sh
export GH_REPO=newrelic/nri-redis
export GH_TAG=v1.6.0
INTEGRATION=nri-redis
VERSION=1.6.0

HAS_LINUX ?= true
HAS_WINDOWS ?= true

OHI_TAR_ARCHS		= "386 amd64 arm arm64"
OHI_TAR				= "${INTEGRATION}_linux_${VERSION}_ARCH.tar.gz"
OHI_WIN_ARCHS		= "386 amd64"
OHI_ZIP				= "${INTEGRATION}-ARCH.${VERSION}.zip"
OHI_MSI				= "${INTEGRATION}-ARCH.${VERSION}.msi"
#OHI_RPM_ARCHS		= "386 1.x86_64 arm arm64"
#OHI_RPM		    = "${INTEGRATION}-${VERSION}-ARCH.rpm"
#OHI_DEB_ARCHS		= "386 amd64 arm arm64"
#OHI_DEB		    = "${INTEGRATION}_${VERSION}-1_ARCH.deb"

default: prepare-secrets mount-s3 mount-s3-check download-ohi-from-gh download-from-gh-check upload-artifacts

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

download-ohi-from-gh:
	@echo "Download packages from github release assets"
ifeq ($(HAS_LINUX),true)
	@source ${DOWNLOAD_SCRIPT} && download ${OHI_TAR} ${OHI_TAR_ARCHS}
endif
ifeq ($(HAS_WINDOWS),true)
	@source ${DOWNLOAD_SCRIPT} && download ${OHI_MSI} ${OHI_WIN_ARCHS}
	@source ${DOWNLOAD_SCRIPT} && download ${OHI_ZIP} ${OHI_WIN_ARCHS}
endif

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



.PHONY: prepare-secrets mount-s3 mount-s3-check download-from-gh download-from-gh-check upload-artifacts