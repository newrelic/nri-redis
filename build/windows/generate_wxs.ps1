param (
	 [int]$major = $(throw "-major is required"),
	 [int]$minor = $(throw "-minor is required"),
	 [int]$patch = $(throw "-patch is required"),
	 [int]$build = 0
)
$integration = "nri-redis"
$integrationName = $integration.Replace("nri-", "")
$executable = "nri-$integrationName.exe"

$projectRootPath = Join-Path -Path $env:GOPATH -ChildPath "src\github.com\newrelic\$integration"

$wix386Path = Join-Path -Path $projectRootPath -ChildPath "build\package\windows\nri-386-installer\Product.wxs"
$wixAmd64Path = Join-Path -Path $projectRootPath -ChildPath "build\package\windows\nri-amd64-installer\Product.wxs"

Function ProcessProductFile($productPath) {
	if ((Test-Path "$productPath.template" -PathType Leaf) -eq $False) {
		Write-Error "$productPath.template not found."
	}
	Copy-Item -Path "$productPath.template" -Destination $productPath -Force

	$product = Get-Content -Path $productPath -Encoding UTF8
	$product = $product -replace "{IntegrationVersion}", "$major.$minor.$patch"
	$product = $product -replace "{Year}", (Get-Date).year
	$product = $product -replace "{IntegrationExe}", $executable
	$product = $product -replace "{IntegrationName}", $integrationName
	Set-Content -Value $product -Path $productPath
}

ProcessProductFile($wix386Path)
ProcessProductFile($wixAmd64Path)