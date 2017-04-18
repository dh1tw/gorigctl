REM Create a release folder
mkdir %GOPATH%\src\github.com\dh1tw\gorigctl\release\

REM copy the needed shared libraries and the binary
%MSYS_PATH%\usr\bin\bash -lc "cp /mingw%MSYS2_BITS%/**/libhamlib-2.dll /c/gopath/src/github.com/dh1tw/gorigctl/release/"
%MSYS_PATH%\usr\bin\bash -lc "cp /mingw%MSYS2_BITS%/**/libgcc_s_dw2-1.dll /c/gopath/src/github.com/dh1tw/gorigctl/release/"
%MSYS_PATH%\usr\bin\bash -lc "cp /mingw%MSYS2_BITS%/**/libwinpthread-1.dll /c/gopath/src/github.com/dh1tw/gorigctl/release/"
REM %MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/dh1tw/gorigctl && ci/release"
%MSYS_PATH%\usr\bin\bash -lc "cp /c/gopath/src/github.com/dh1tw/gorigctl/gorigctl.exe /c/gopath/src/github.com/dh1tw/gorigctl/release"

REM zip everything
%MSYS_PATH%\usr\bin\bash -lc "cd /c/gopath/src/github.com/dh1tw/gorigctl/release && 7z a -tzip gorigctl-v$APPVEYOR_REPO_TAG_NAME-$GOOS-$GOARCH.zip *"

REM copy it into the build folder
xcopy %GOPATH%\src\github.com\dh1tw\gorigctl\release\gorigctl-v%APPVEYOR_REPO_TAG_NAME%-%GOOS%-%GOARCH%.zip %APPVEYOR_BUILD_FOLDER%\ /e /i > nul