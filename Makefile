
all:
	go build

install: build_it	
	cp check-json-syntax ~/bin

# DIFF=cp
# DIFF=cat
DIFF=diff

test: build_it test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test13 test14 test15
	@echo PASS

build_it:
	go build

test1: build_it
	./check-json-syntax test/t1.json >out/t1.out
	$(DIFF) out/t1.out ref/t1.ref
	@echo PASS

test2: build_it
	./check-json-syntax test/t2.json >out/t2.out
	$(DIFF) out/t2.out ref/t2.ref
	@echo PASS

test3: build_it
	./check-json-syntax test/t3.json >out/t3.out
	$(DIFF) out/t3.out ref/t3.ref
	@echo PASS

# check array
test4: build_it
	./check-json-syntax test/t4.json >out/t4.out
	$(DIFF) out/t4.out ref/t4.ref
	@echo PASS

test5: build_it
	./check-json-syntax test/t5.json >out/t5.out
	$(DIFF) out/t5.out ref/t5.ref
	@echo PASS

test6: build_it
	./check-json-syntax test/t6.json >out/t6.out
	$(DIFF) out/t6.out ref/t6.ref
	@echo PASS

test7: build_it
	./check-json-syntax test/t7.json >out/t7.out
	$(DIFF) out/t7.out ref/t7.ref
	@echo PASS

test8: build_it
	./check-json-syntax test/t8.json >out/t8.out
	$(DIFF) out/t8.out ref/t8.ref
	@echo PASS

test9: build_it
	./check-json-syntax test/t9.json >out/t9.out
	$(DIFF) out/t9.out ref/t9.ref
	@echo PASS

test10: build_it
	./check-json-syntax test/t10.json >out/t10.out
	$(DIFF) out/t10.out ref/t10.ref
	@echo PASS

test11: build_it
	./check-json-syntax test/t11.json >out/t11.out
	$(DIFF) out/t11.out ref/t11.ref
	@echo PASS

test12: build_it
	./check-json-syntax test/t12.json >out/t12.out
	$(DIFF) out/t12.out ref/t12.ref
	@echo PASS

test13: build_it
	./check-json-syntax test/t13.json >out/t13.out
	$(DIFF) out/t13.out ref/t13.ref
	@echo PASS

test14: build_it
	./check-json-syntax -l test/t14.json >out/t14.out
	$(DIFF) out/t14.out ref/t14.ref
	@echo PASS

test15: build_it
	./check-json-syntax -p test/t15.json >out/t15.out
	$(DIFF) out/t15.out ref/t15.ref
	@echo PASS

