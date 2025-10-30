window.BENCHMARK_DATA = {
  "lastUpdate": 1761856570501,
  "repoUrl": "https://github.com/SPANDigital/cel2sql",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "richard.wooding@gmail.com",
            "name": "Richard Wooding",
            "username": "richardwooding"
          },
          "committer": {
            "email": "richard.wooding@gmail.com",
            "name": "Richard Wooding",
            "username": "richardwooding"
          },
          "distinct": true,
          "id": "429a6985a064c478a1487aa0e79ea00293405699",
          "message": "fix: Add write permissions to benchmark job for gh-pages push\n\nThe benchmark job needs 'contents: write' permission to push\nbenchmark results to the gh-pages branch.\n\nFixes: remote: Permission to SPANDigital/cel2sql.git denied to github-actions[bot]\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-Authored-By: Claude <noreply@anthropic.com>",
          "timestamp": "2025-10-30T22:33:53+02:00",
          "tree_id": "b8d6dde609e6fc5d5ce47ae9e09d46e60b5b4744",
          "url": "https://github.com/SPANDigital/cel2sql/commit/429a6985a064c478a1487aa0e79ea00293405699"
        },
        "date": 1761856569873,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2001,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "585944 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2001,
            "unit": "ns/op",
            "extra": "585944 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "585944 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "585944 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2022,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "582434 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2022,
            "unit": "ns/op",
            "extra": "582434 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "582434 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "582434 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2161,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "535117 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2161,
            "unit": "ns/op",
            "extra": "535117 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "535117 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "535117 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1390,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "852231 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1390,
            "unit": "ns/op",
            "extra": "852231 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "852231 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "852231 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 2991,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "386310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 2991,
            "unit": "ns/op",
            "extra": "386310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "386310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "386310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3664,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "325105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3664,
            "unit": "ns/op",
            "extra": "325105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "325105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "325105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2760,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "430932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2760,
            "unit": "ns/op",
            "extra": "430932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "430932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "430932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2748,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "431342 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2748,
            "unit": "ns/op",
            "extra": "431342 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "431342 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "431342 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3519,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "378042 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3519,
            "unit": "ns/op",
            "extra": "378042 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "378042 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "378042 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3121,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "376045 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3121,
            "unit": "ns/op",
            "extra": "376045 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "376045 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "376045 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2718,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "445938 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2718,
            "unit": "ns/op",
            "extra": "445938 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "445938 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "445938 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3067,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "344259 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3067,
            "unit": "ns/op",
            "extra": "344259 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "344259 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "344259 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7385,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "161545 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7385,
            "unit": "ns/op",
            "extra": "161545 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "161545 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "161545 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5247,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "222633 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5247,
            "unit": "ns/op",
            "extra": "222633 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "222633 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "222633 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7607,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "156939 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7607,
            "unit": "ns/op",
            "extra": "156939 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "156939 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "156939 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5949,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "200211 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5949,
            "unit": "ns/op",
            "extra": "200211 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "200211 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "200211 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7670,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "155875 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7670,
            "unit": "ns/op",
            "extra": "155875 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "155875 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "155875 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6938,
            "unit": "ns/op\t    5674 B/op\t      88 allocs/op",
            "extra": "174588 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6938,
            "unit": "ns/op",
            "extra": "174588 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5674,
            "unit": "B/op",
            "extra": "174588 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "174588 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7488,
            "unit": "ns/op\t    5841 B/op\t     105 allocs/op",
            "extra": "159710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7488,
            "unit": "ns/op",
            "extra": "159710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5841,
            "unit": "B/op",
            "extra": "159710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "159710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 6071,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "203359 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 6071,
            "unit": "ns/op",
            "extra": "203359 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "203359 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "203359 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3519,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "338857 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3519,
            "unit": "ns/op",
            "extra": "338857 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "338857 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "338857 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4122,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "287269 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4122,
            "unit": "ns/op",
            "extra": "287269 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "287269 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "287269 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2802,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "425244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2802,
            "unit": "ns/op",
            "extra": "425244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "425244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "425244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3268,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "367340 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3268,
            "unit": "ns/op",
            "extra": "367340 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "367340 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "367340 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3481,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "345146 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3481,
            "unit": "ns/op",
            "extra": "345146 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "345146 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "345146 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7444,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "158142 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7444,
            "unit": "ns/op",
            "extra": "158142 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "158142 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "158142 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8251,
            "unit": "ns/op\t    5974 B/op\t      75 allocs/op",
            "extra": "143248 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8251,
            "unit": "ns/op",
            "extra": "143248 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5974,
            "unit": "B/op",
            "extra": "143248 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "143248 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7930,
            "unit": "ns/op\t    5979 B/op\t      75 allocs/op",
            "extra": "150168 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7930,
            "unit": "ns/op",
            "extra": "150168 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5979,
            "unit": "B/op",
            "extra": "150168 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "150168 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8376,
            "unit": "ns/op\t    5999 B/op\t      75 allocs/op",
            "extra": "143428 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8376,
            "unit": "ns/op",
            "extra": "143428 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5999,
            "unit": "B/op",
            "extra": "143428 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "143428 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8137,
            "unit": "ns/op\t    5992 B/op\t      76 allocs/op",
            "extra": "147127 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8137,
            "unit": "ns/op",
            "extra": "147127 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5992,
            "unit": "B/op",
            "extra": "147127 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147127 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7952,
            "unit": "ns/op\t    5990 B/op\t      76 allocs/op",
            "extra": "147596 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7952,
            "unit": "ns/op",
            "extra": "147596 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5990,
            "unit": "B/op",
            "extra": "147596 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147596 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8188,
            "unit": "ns/op\t    5988 B/op\t      76 allocs/op",
            "extra": "144363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8188,
            "unit": "ns/op",
            "extra": "144363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5988,
            "unit": "B/op",
            "extra": "144363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "144363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9683,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "122628 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9683,
            "unit": "ns/op",
            "extra": "122628 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "122628 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "122628 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 19043,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "62842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 19043,
            "unit": "ns/op",
            "extra": "62842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "62842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "62842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2095,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "567925 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2095,
            "unit": "ns/op",
            "extra": "567925 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "567925 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "567925 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8755,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "138310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8755,
            "unit": "ns/op",
            "extra": "138310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "138310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "138310 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5410,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "218084 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5410,
            "unit": "ns/op",
            "extra": "218084 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "218084 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "218084 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42492,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28219 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42492,
            "unit": "ns/op",
            "extra": "28219 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28219 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28219 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39185,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30717 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39185,
            "unit": "ns/op",
            "extra": "30717 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30717 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30717 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3764,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "312460 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3764,
            "unit": "ns/op",
            "extra": "312460 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "312460 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "312460 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1804,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "652843 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1804,
            "unit": "ns/op",
            "extra": "652843 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "652843 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "652843 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1864,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "617817 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1864,
            "unit": "ns/op",
            "extra": "617817 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "617817 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "617817 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2238,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "540212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2238,
            "unit": "ns/op",
            "extra": "540212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "540212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "540212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2275,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "522742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2275,
            "unit": "ns/op",
            "extra": "522742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "522742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "522742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2181,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "543770 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2181,
            "unit": "ns/op",
            "extra": "543770 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "543770 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "543770 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3206,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "381044 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3206,
            "unit": "ns/op",
            "extra": "381044 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "381044 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "381044 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6164,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "193136 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6164,
            "unit": "ns/op",
            "extra": "193136 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "193136 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193136 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6225,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "193659 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6225,
            "unit": "ns/op",
            "extra": "193659 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "193659 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193659 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6139,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "194716 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6139,
            "unit": "ns/op",
            "extra": "194716 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "194716 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "194716 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12207,
            "unit": "ns/op\t    8326 B/op\t     117 allocs/op",
            "extra": "99128 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12207,
            "unit": "ns/op",
            "extra": "99128 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8326,
            "unit": "B/op",
            "extra": "99128 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "99128 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6858,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "172376 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6858,
            "unit": "ns/op",
            "extra": "172376 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "172376 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "172376 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22431,
            "unit": "ns/op\t   13925 B/op\t     217 allocs/op",
            "extra": "53390 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22431,
            "unit": "ns/op",
            "extra": "53390 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13925,
            "unit": "B/op",
            "extra": "53390 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "53390 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2788,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "415474 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2788,
            "unit": "ns/op",
            "extra": "415474 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "415474 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "415474 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2795,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "433204 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2795,
            "unit": "ns/op",
            "extra": "433204 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "433204 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "433204 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2797,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "418094 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2797,
            "unit": "ns/op",
            "extra": "418094 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "418094 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "418094 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2796,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "425187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2796,
            "unit": "ns/op",
            "extra": "425187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "425187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "425187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2810,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "414604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2810,
            "unit": "ns/op",
            "extra": "414604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "414604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "414604 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 607338,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1956 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 607338,
            "unit": "ns/op",
            "extra": "1956 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1956 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1956 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 500656,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 500656,
            "unit": "ns/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2377 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 620.9,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1965517 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 620.9,
            "unit": "ns/op",
            "extra": "1965517 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1965517 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1965517 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 623.4,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1940857 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 623.4,
            "unit": "ns/op",
            "extra": "1940857 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1940857 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1940857 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 615.8,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1947296 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 615.8,
            "unit": "ns/op",
            "extra": "1947296 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1947296 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1947296 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 142,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8451036 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 142,
            "unit": "ns/op",
            "extra": "8451036 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8451036 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8451036 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 8910,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "152116 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 8910,
            "unit": "ns/op",
            "extra": "152116 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "152116 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "152116 times\n4 procs"
          }
        ]
      }
    ]
  }
}