rm -rf _builds _steps _projects _cache _temp
wercker build --git-domain github.com --git-owner cloudnativego --git-repository wof-mapservice
rm -rf _builds _steps _projects _cache _temp 
