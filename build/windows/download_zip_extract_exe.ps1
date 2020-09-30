param (
    [string]$INTEGRATION="none",
    [string]$ARCH="amd64",
    [string]$TAG="v0.0.0",
    [string]$REPO_FULL_NAME="none"
)
$VERSION=${TAG}.substring(1)
$zip_name="nri-${INTEGRATION}_windows_${VERSION}_${ARCH}.zip"

write-host "===> Downloading & extracting .exe from https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${zip_name}"

Invoke-WebRequest "https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${zip_name}" -OutFile ".\dist\${zip_name}"
expand-archive -path '.\dist\${zip_name}' -destinationpath .\dist\
Copy-Item -Path ".\dist\New Relic\newrelic-infra\newrelic-integrations\bin\nri-${INTEGRATION}.exe" -Destination ".\dist\" -Force