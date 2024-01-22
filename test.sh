# 현재 디렉토리부터 시작하여 모든 하위 디렉토리를 탐색
for dir in $(find . -type d);
do
    # Go 파일과 테스트 파일이 모두 있는지 확인
    if ls ${dir}/*.go &> /dev/null && ls ${dir}/*_test.go &> /dev/null; then
        echo "Testing ${dir}"
        go test -cover ${dir}
    fi
done
