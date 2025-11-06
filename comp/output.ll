@format_read_int = global [3 x i8] c"%d\00"
@format_write_int = global [4 x i8] c"%d\0A\00"

declare i32 @scanf(i8* %format, ...)

declare i32 @printf(i8* %format, ...)

define i32 @add(i32 %a, i32 %b) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = add i32 %a, %b
	store i32 %2, i32* %1
	%3 = load i32, i32* %1
	ret i32 %3
}

define i32 @minus(i32 %a, i32 %b) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = sub i32 %a, %b
	store i32 %2, i32* %1
	%3 = load i32, i32* %1
	ret i32 %3
}

define i32 @times(i32 %a, i32 %b) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = mul i32 %a, %b
	store i32 %2, i32* %1
	%3 = load i32, i32* %1
	ret i32 %3
}

define i32 @div(i32 %a, i32 %b) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = sdiv i32 %a, %b
	store i32 %2, i32* %1
	%3 = load i32, i32* %1
	ret i32 %3
}

define i32 @main() {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = getelementptr [3 x i8], [3 x i8]* @format_read_int, i32 0, i32 0
	%3 = call i32 (i8*, ...) @scanf(i8* %2, i32* %1)
	%4 = alloca i32
	store i32 0, i32* %4
	%5 = call i32 (i8*, ...) @scanf(i8* %2, i32* %4)
	%6 = alloca i32
	store i32 0, i32* %6
	%7 = load i32, i32* %1
	%8 = load i32, i32* %4
	%9 = call i32 @add(i32 %7, i32 %8)
	store i32 %9, i32* %6
	%10 = load i32, i32* %6
	%11 = getelementptr [4 x i8], [4 x i8]* @format_write_int, i32 0, i32 0
	%12 = call i32 (i8*, ...) @printf(i8* %11, i32 %10)
	%13 = load i32, i32* %1
	%14 = load i32, i32* %4
	%15 = call i32 @minus(i32 %13, i32 %14)
	store i32 %15, i32* %6
	%16 = load i32, i32* %6
	%17 = call i32 (i8*, ...) @printf(i8* %11, i32 %16)
	%18 = load i32, i32* %1
	%19 = load i32, i32* %4
	%20 = call i32 @times(i32 %18, i32 %19)
	store i32 %20, i32* %6
	%21 = load i32, i32* %6
	%22 = call i32 (i8*, ...) @printf(i8* %11, i32 %21)
	%23 = load i32, i32* %1
	%24 = load i32, i32* %4
	%25 = call i32 @div(i32 %23, i32 %24)
	store i32 %25, i32* %6
	%26 = load i32, i32* %6
	%27 = call i32 (i8*, ...) @printf(i8* %11, i32 %26)
	%28 = load i32, i32* %1
	ret i32 %28
}
