param (
    [string]$INTEGRATION="none",
    [string]$ARCH="amd64",
    [string]$TAG="v0.0.0",
    [string]$REPO_FULL_NAME="none"
)
$VERSION=${TAG}.substring(1)
$exe_folder="nri-${INTEGRATION}_windows_${ARCH}"
$zip_name="nri-${INTEGRATION}_windows_${VERSION}_${ARCH}.zip"

$zip_url="https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${zip_name}"
write-host "===> Downloading & extracting .exe from ${zip_url}"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
Invoke-WebRequest "${zip_url}" -OutFile ".\dist\${zip_name}"
write-host "===> Expanding"
expand-archive -path "dist\${zip_name}" -destinationpath "dist\${exe_folder}\"
write-host "===> copying"
Copy-Item -Path ".\dist\${exe_folder}\New Relic\newrelic-infra\newrelic-integrations\bin\nri-${INTEGRATION}.exe" -Destination "dist\${exe_folder}\" -Force