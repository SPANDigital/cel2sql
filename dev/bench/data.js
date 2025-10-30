window.BENCHMARK_DATA = {
  "lastUpdate": 1761863511471,
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
      },
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
          "id": "7b770aa8931368f229a7b3d67543c175cec15c4d",
          "message": "docs: Add benchmark tracking badge to README\n\n- Add badge linking to GitHub Pages benchmark results\n- Update features list to mention continuous benchmark monitoring\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-Authored-By: Claude <noreply@anthropic.com>",
          "timestamp": "2025-10-30T22:42:23+02:00",
          "tree_id": "fb48011775cb28b55fbb5019a9c68d89427da487",
          "url": "https://github.com/SPANDigital/cel2sql/commit/7b770aa8931368f229a7b3d67543c175cec15c4d"
        },
        "date": 1761857077161,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 1985,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "595809 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 1985,
            "unit": "ns/op",
            "extra": "595809 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "595809 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "595809 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2022,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "584023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2022,
            "unit": "ns/op",
            "extra": "584023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "584023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "584023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2447,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "542864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2447,
            "unit": "ns/op",
            "extra": "542864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "542864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "542864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1372,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "867314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1372,
            "unit": "ns/op",
            "extra": "867314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "867314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "867314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 2987,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "392317 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 2987,
            "unit": "ns/op",
            "extra": "392317 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "392317 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "392317 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3628,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "329007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3628,
            "unit": "ns/op",
            "extra": "329007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "329007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "329007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2754,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "431083 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2754,
            "unit": "ns/op",
            "extra": "431083 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "431083 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "431083 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2738,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "433242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2738,
            "unit": "ns/op",
            "extra": "433242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "433242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "433242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3110,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "379396 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3110,
            "unit": "ns/op",
            "extra": "379396 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "379396 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "379396 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3106,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "386620 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3106,
            "unit": "ns/op",
            "extra": "386620 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "386620 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "386620 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2687,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "444709 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2687,
            "unit": "ns/op",
            "extra": "444709 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "444709 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "444709 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3082,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "385614 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3082,
            "unit": "ns/op",
            "extra": "385614 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "385614 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "385614 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7361,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "160465 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7361,
            "unit": "ns/op",
            "extra": "160465 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "160465 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "160465 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5282,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "226932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5282,
            "unit": "ns/op",
            "extra": "226932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "226932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "226932 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7572,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7572,
            "unit": "ns/op",
            "extra": "158703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5928,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "194832 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5928,
            "unit": "ns/op",
            "extra": "194832 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "194832 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "194832 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7563,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "158665 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7563,
            "unit": "ns/op",
            "extra": "158665 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "158665 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "158665 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6957,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "167608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6957,
            "unit": "ns/op",
            "extra": "167608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "167608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "167608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7354,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "164214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7354,
            "unit": "ns/op",
            "extra": "164214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "164214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "164214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 6042,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "200778 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 6042,
            "unit": "ns/op",
            "extra": "200778 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "200778 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "200778 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3465,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "342462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3465,
            "unit": "ns/op",
            "extra": "342462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "342462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "342462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4103,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "288054 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4103,
            "unit": "ns/op",
            "extra": "288054 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "288054 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "288054 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2780,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "422498 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2780,
            "unit": "ns/op",
            "extra": "422498 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "422498 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "422498 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3230,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "367922 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3230,
            "unit": "ns/op",
            "extra": "367922 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "367922 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "367922 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3465,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "294710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3465,
            "unit": "ns/op",
            "extra": "294710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "294710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "294710 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7478,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "158943 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7478,
            "unit": "ns/op",
            "extra": "158943 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "158943 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "158943 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8263,
            "unit": "ns/op\t    5972 B/op\t      75 allocs/op",
            "extra": "144481 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8263,
            "unit": "ns/op",
            "extra": "144481 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5972,
            "unit": "B/op",
            "extra": "144481 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "144481 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7879,
            "unit": "ns/op\t    5980 B/op\t      75 allocs/op",
            "extra": "152199 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7879,
            "unit": "ns/op",
            "extra": "152199 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5980,
            "unit": "B/op",
            "extra": "152199 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "152199 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8398,
            "unit": "ns/op\t    5988 B/op\t      75 allocs/op",
            "extra": "141898 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8398,
            "unit": "ns/op",
            "extra": "141898 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5988,
            "unit": "B/op",
            "extra": "141898 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "141898 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8079,
            "unit": "ns/op\t    5986 B/op\t      76 allocs/op",
            "extra": "149296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8079,
            "unit": "ns/op",
            "extra": "149296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5986,
            "unit": "B/op",
            "extra": "149296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "149296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7940,
            "unit": "ns/op\t    5990 B/op\t      76 allocs/op",
            "extra": "151904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7940,
            "unit": "ns/op",
            "extra": "151904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5990,
            "unit": "B/op",
            "extra": "151904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "151904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8152,
            "unit": "ns/op\t    5986 B/op\t      76 allocs/op",
            "extra": "147244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8152,
            "unit": "ns/op",
            "extra": "147244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5986,
            "unit": "B/op",
            "extra": "147244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9722,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "123462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9722,
            "unit": "ns/op",
            "extra": "123462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "123462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "123462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 18877,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "63276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 18877,
            "unit": "ns/op",
            "extra": "63276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "63276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "63276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2092,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "560176 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2092,
            "unit": "ns/op",
            "extra": "560176 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "560176 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "560176 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8570,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "137053 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8570,
            "unit": "ns/op",
            "extra": "137053 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "137053 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "137053 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5437,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "219332 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5437,
            "unit": "ns/op",
            "extra": "219332 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "219332 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "219332 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42134,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28101 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42134,
            "unit": "ns/op",
            "extra": "28101 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28101 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28101 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 38985,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30727 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 38985,
            "unit": "ns/op",
            "extra": "30727 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30727 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30727 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3767,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "313335 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3767,
            "unit": "ns/op",
            "extra": "313335 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "313335 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "313335 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1800,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "657025 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1800,
            "unit": "ns/op",
            "extra": "657025 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "657025 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "657025 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1827,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "656006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1827,
            "unit": "ns/op",
            "extra": "656006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "656006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "656006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2212,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "535654 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2212,
            "unit": "ns/op",
            "extra": "535654 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "535654 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "535654 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2285,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "519414 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2285,
            "unit": "ns/op",
            "extra": "519414 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "519414 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "519414 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2163,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "538759 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2163,
            "unit": "ns/op",
            "extra": "538759 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "538759 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "538759 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3223,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "372406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3223,
            "unit": "ns/op",
            "extra": "372406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "372406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "372406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6114,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "196578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6114,
            "unit": "ns/op",
            "extra": "196578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "196578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "196578 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6149,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "192970 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6149,
            "unit": "ns/op",
            "extra": "192970 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "192970 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "192970 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6068,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "195728 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6068,
            "unit": "ns/op",
            "extra": "195728 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "195728 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "195728 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12372,
            "unit": "ns/op\t    8321 B/op\t     117 allocs/op",
            "extra": "99049 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12372,
            "unit": "ns/op",
            "extra": "99049 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8321,
            "unit": "B/op",
            "extra": "99049 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "99049 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6861,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "175470 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6861,
            "unit": "ns/op",
            "extra": "175470 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "175470 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "175470 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22445,
            "unit": "ns/op\t   13935 B/op\t     217 allocs/op",
            "extra": "54020 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22445,
            "unit": "ns/op",
            "extra": "54020 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13935,
            "unit": "B/op",
            "extra": "54020 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "54020 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2801,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "400708 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2801,
            "unit": "ns/op",
            "extra": "400708 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "400708 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "400708 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2789,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "424149 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2789,
            "unit": "ns/op",
            "extra": "424149 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "424149 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "424149 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2788,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "421804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2788,
            "unit": "ns/op",
            "extra": "421804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "421804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "421804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2782,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "423254 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2782,
            "unit": "ns/op",
            "extra": "423254 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "423254 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "423254 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2796,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "416102 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2796,
            "unit": "ns/op",
            "extra": "416102 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "416102 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "416102 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 604803,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1945 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 604803,
            "unit": "ns/op",
            "extra": "1945 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1945 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1945 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 498792,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 498792,
            "unit": "ns/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 673.7,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1957358 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 673.7,
            "unit": "ns/op",
            "extra": "1957358 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1957358 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1957358 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 624.5,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1885142 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 624.5,
            "unit": "ns/op",
            "extra": "1885142 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1885142 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1885142 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 621.2,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1945982 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 621.2,
            "unit": "ns/op",
            "extra": "1945982 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1945982 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1945982 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 146.5,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8283364 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 146.5,
            "unit": "ns/op",
            "extra": "8283364 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8283364 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8283364 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 7872,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "153679 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 7872,
            "unit": "ns/op",
            "extra": "153679 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "153679 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "153679 times\n4 procs"
          }
        ]
      },
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
          "id": "08684a33eea8b027adf7453910f0de7c1052540a",
          "message": "chore: Add benchmark comparison files to .gitignore\n\n- Add bench-old.txt and bench-new.txt to .gitignore\n- These files are generated by make bench-compare for local benchmark comparisons\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-Authored-By: Claude <noreply@anthropic.com>",
          "timestamp": "2025-10-30T22:43:30+02:00",
          "tree_id": "76d4eedebb883d8f5ea2f87432a5387eba065af4",
          "url": "https://github.com/SPANDigital/cel2sql/commit/08684a33eea8b027adf7453910f0de7c1052540a"
        },
        "date": 1761857149634,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 1997,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "590660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 1997,
            "unit": "ns/op",
            "extra": "590660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "590660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "590660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2284,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "571636 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2284,
            "unit": "ns/op",
            "extra": "571636 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "571636 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "571636 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2157,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "545029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2157,
            "unit": "ns/op",
            "extra": "545029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "545029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "545029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1371,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "811092 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1371,
            "unit": "ns/op",
            "extra": "811092 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "811092 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "811092 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 2994,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "392872 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 2994,
            "unit": "ns/op",
            "extra": "392872 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "392872 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "392872 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3643,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "332726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3643,
            "unit": "ns/op",
            "extra": "332726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "332726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "332726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2743,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "434494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2743,
            "unit": "ns/op",
            "extra": "434494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "434494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "434494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2739,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "436657 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2739,
            "unit": "ns/op",
            "extra": "436657 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "436657 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "436657 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3118,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "374974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3118,
            "unit": "ns/op",
            "extra": "374974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "374974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "374974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3128,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "384006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3128,
            "unit": "ns/op",
            "extra": "384006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "384006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "384006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2694,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "437578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2694,
            "unit": "ns/op",
            "extra": "437578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "437578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "437578 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3052,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "390680 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3052,
            "unit": "ns/op",
            "extra": "390680 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "390680 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "390680 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7448,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "159984 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7448,
            "unit": "ns/op",
            "extra": "159984 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "159984 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "159984 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5260,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "229226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5260,
            "unit": "ns/op",
            "extra": "229226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "229226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "229226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7579,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158854 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7579,
            "unit": "ns/op",
            "extra": "158854 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158854 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158854 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5936,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "203436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5936,
            "unit": "ns/op",
            "extra": "203436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "203436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "203436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7671,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "156909 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7671,
            "unit": "ns/op",
            "extra": "156909 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "156909 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "156909 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6966,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "171747 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6966,
            "unit": "ns/op",
            "extra": "171747 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "171747 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "171747 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7484,
            "unit": "ns/op\t    5841 B/op\t     105 allocs/op",
            "extra": "162078 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7484,
            "unit": "ns/op",
            "extra": "162078 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5841,
            "unit": "B/op",
            "extra": "162078 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "162078 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 5939,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "200034 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 5939,
            "unit": "ns/op",
            "extra": "200034 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "200034 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "200034 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3473,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "341029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3473,
            "unit": "ns/op",
            "extra": "341029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "341029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "341029 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4098,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "292165 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4098,
            "unit": "ns/op",
            "extra": "292165 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "292165 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "292165 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2892,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "415130 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2892,
            "unit": "ns/op",
            "extra": "415130 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "415130 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "415130 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3255,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "366896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3255,
            "unit": "ns/op",
            "extra": "366896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "366896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "366896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3451,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "340530 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3451,
            "unit": "ns/op",
            "extra": "340530 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "340530 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "340530 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7564,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "160542 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7564,
            "unit": "ns/op",
            "extra": "160542 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "160542 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "160542 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8271,
            "unit": "ns/op\t    5965 B/op\t      75 allocs/op",
            "extra": "142876 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8271,
            "unit": "ns/op",
            "extra": "142876 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5965,
            "unit": "B/op",
            "extra": "142876 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142876 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7944,
            "unit": "ns/op\t    5971 B/op\t      75 allocs/op",
            "extra": "151240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7944,
            "unit": "ns/op",
            "extra": "151240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5971,
            "unit": "B/op",
            "extra": "151240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "151240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8409,
            "unit": "ns/op\t    5995 B/op\t      75 allocs/op",
            "extra": "143671 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8409,
            "unit": "ns/op",
            "extra": "143671 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5995,
            "unit": "B/op",
            "extra": "143671 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "143671 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8115,
            "unit": "ns/op\t    5991 B/op\t      76 allocs/op",
            "extra": "143827 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8115,
            "unit": "ns/op",
            "extra": "143827 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5991,
            "unit": "B/op",
            "extra": "143827 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "143827 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7950,
            "unit": "ns/op\t    5990 B/op\t      76 allocs/op",
            "extra": "151550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7950,
            "unit": "ns/op",
            "extra": "151550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5990,
            "unit": "B/op",
            "extra": "151550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "151550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8197,
            "unit": "ns/op\t    5991 B/op\t      76 allocs/op",
            "extra": "146082 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8197,
            "unit": "ns/op",
            "extra": "146082 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5991,
            "unit": "B/op",
            "extra": "146082 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "146082 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9733,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "123974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9733,
            "unit": "ns/op",
            "extra": "123974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "123974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "123974 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 18907,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "62630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 18907,
            "unit": "ns/op",
            "extra": "62630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "62630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "62630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2097,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "569424 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2097,
            "unit": "ns/op",
            "extra": "569424 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "569424 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "569424 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8679,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "138764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8679,
            "unit": "ns/op",
            "extra": "138764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "138764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "138764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5427,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "218455 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5427,
            "unit": "ns/op",
            "extra": "218455 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "218455 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "218455 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42447,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28394 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42447,
            "unit": "ns/op",
            "extra": "28394 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28394 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28394 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39324,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30858 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39324,
            "unit": "ns/op",
            "extra": "30858 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30858 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30858 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3761,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "316483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3761,
            "unit": "ns/op",
            "extra": "316483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "316483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "316483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1809,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "661388 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1809,
            "unit": "ns/op",
            "extra": "661388 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "661388 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "661388 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1834,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "639907 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1834,
            "unit": "ns/op",
            "extra": "639907 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "639907 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "639907 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2232,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "533698 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2232,
            "unit": "ns/op",
            "extra": "533698 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "533698 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "533698 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2268,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "500043 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2268,
            "unit": "ns/op",
            "extra": "500043 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "500043 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "500043 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2206,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "547795 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2206,
            "unit": "ns/op",
            "extra": "547795 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "547795 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "547795 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3173,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "380685 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3173,
            "unit": "ns/op",
            "extra": "380685 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "380685 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "380685 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6184,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "194014 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6184,
            "unit": "ns/op",
            "extra": "194014 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "194014 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "194014 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6264,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "192318 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6264,
            "unit": "ns/op",
            "extra": "192318 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "192318 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "192318 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6150,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "195536 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6150,
            "unit": "ns/op",
            "extra": "195536 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "195536 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "195536 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12162,
            "unit": "ns/op\t    8321 B/op\t     117 allocs/op",
            "extra": "98539 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12162,
            "unit": "ns/op",
            "extra": "98539 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8321,
            "unit": "B/op",
            "extra": "98539 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "98539 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6907,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "173287 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6907,
            "unit": "ns/op",
            "extra": "173287 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "173287 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "173287 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22507,
            "unit": "ns/op\t   13933 B/op\t     217 allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22507,
            "unit": "ns/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13933,
            "unit": "B/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2788,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "423421 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2788,
            "unit": "ns/op",
            "extra": "423421 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "423421 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "423421 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2786,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "427135 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2786,
            "unit": "ns/op",
            "extra": "427135 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "427135 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "427135 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2781,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "417066 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2781,
            "unit": "ns/op",
            "extra": "417066 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "417066 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "417066 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2786,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "424971 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2786,
            "unit": "ns/op",
            "extra": "424971 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "424971 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "424971 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2801,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "418350 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2801,
            "unit": "ns/op",
            "extra": "418350 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "418350 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "418350 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 604748,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1962 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 604748,
            "unit": "ns/op",
            "extra": "1962 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1962 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1962 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 494662,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 494662,
            "unit": "ns/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 606.5,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1975221 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 606.5,
            "unit": "ns/op",
            "extra": "1975221 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1975221 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1975221 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 615.1,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1948056 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 615.1,
            "unit": "ns/op",
            "extra": "1948056 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1948056 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1948056 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 614.3,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1948567 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 614.3,
            "unit": "ns/op",
            "extra": "1948567 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1948567 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1948567 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 143.5,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8437436 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 143.5,
            "unit": "ns/op",
            "extra": "8437436 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8437436 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8437436 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 8568,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "152409 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 8568,
            "unit": "ns/op",
            "extra": "152409 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "152409 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "152409 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "SPANDigital",
            "username": "SPANDigital"
          },
          "committer": {
            "name": "SPANDigital",
            "username": "SPANDigital"
          },
          "id": "6461402b0114acc19b5dee6b40be5d66d6ab4374",
          "message": "fix: Add byte array length validation to prevent resource exhaustion (fixes #36)",
          "timestamp": "2025-10-30T20:43:34Z",
          "url": "https://github.com/SPANDigital/cel2sql/pull/82/commits/6461402b0114acc19b5dee6b40be5d66d6ab4374"
        },
        "date": 1761861026475,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkByteArrayConversion",
            "value": 8135,
            "unit": "ns/op\t   14708 B/op\t      30 allocs/op",
            "extra": "150223 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - ns/op",
            "value": 8135,
            "unit": "ns/op",
            "extra": "150223 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - B/op",
            "value": 14708,
            "unit": "B/op",
            "extra": "150223 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "150223 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2016,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "576992 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2016,
            "unit": "ns/op",
            "extra": "576992 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "576992 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "576992 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2033,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "585195 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2033,
            "unit": "ns/op",
            "extra": "585195 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "585195 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "585195 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2167,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "538766 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2167,
            "unit": "ns/op",
            "extra": "538766 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "538766 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "538766 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1391,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "846890 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1391,
            "unit": "ns/op",
            "extra": "846890 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "846890 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "846890 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 3018,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "386834 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 3018,
            "unit": "ns/op",
            "extra": "386834 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "386834 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "386834 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3648,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "319407 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3648,
            "unit": "ns/op",
            "extra": "319407 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "319407 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "319407 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2930,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "435799 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2930,
            "unit": "ns/op",
            "extra": "435799 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "435799 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "435799 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2940,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "376164 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2940,
            "unit": "ns/op",
            "extra": "376164 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "376164 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "376164 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3144,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "379276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3144,
            "unit": "ns/op",
            "extra": "379276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "379276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "379276 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3163,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "370441 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3163,
            "unit": "ns/op",
            "extra": "370441 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "370441 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "370441 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2715,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "438846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2715,
            "unit": "ns/op",
            "extra": "438846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "438846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "438846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3099,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "371856 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3099,
            "unit": "ns/op",
            "extra": "371856 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "371856 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "371856 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7397,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "162139 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7397,
            "unit": "ns/op",
            "extra": "162139 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "162139 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "162139 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5283,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "228652 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5283,
            "unit": "ns/op",
            "extra": "228652 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "228652 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "228652 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7589,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7589,
            "unit": "ns/op",
            "extra": "158968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5959,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "203240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5959,
            "unit": "ns/op",
            "extra": "203240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "203240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "203240 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7671,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "156734 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7671,
            "unit": "ns/op",
            "extra": "156734 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "156734 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "156734 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 7064,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "171213 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 7064,
            "unit": "ns/op",
            "extra": "171213 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "171213 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "171213 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7352,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "158118 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7352,
            "unit": "ns/op",
            "extra": "158118 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "158118 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "158118 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 5926,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "201729 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 5926,
            "unit": "ns/op",
            "extra": "201729 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "201729 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "201729 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3532,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "335845 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3532,
            "unit": "ns/op",
            "extra": "335845 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "335845 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "335845 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4147,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "283390 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4147,
            "unit": "ns/op",
            "extra": "283390 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "283390 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "283390 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2817,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "427378 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2817,
            "unit": "ns/op",
            "extra": "427378 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "427378 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "427378 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3250,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "364087 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3250,
            "unit": "ns/op",
            "extra": "364087 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "364087 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "364087 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3499,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "338486 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3499,
            "unit": "ns/op",
            "extra": "338486 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "338486 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "338486 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7541,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "161700 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7541,
            "unit": "ns/op",
            "extra": "161700 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "161700 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "161700 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8349,
            "unit": "ns/op\t    5977 B/op\t      75 allocs/op",
            "extra": "142446 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8349,
            "unit": "ns/op",
            "extra": "142446 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5977,
            "unit": "B/op",
            "extra": "142446 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142446 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 8021,
            "unit": "ns/op\t    5977 B/op\t      75 allocs/op",
            "extra": "144289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 8021,
            "unit": "ns/op",
            "extra": "144289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5977,
            "unit": "B/op",
            "extra": "144289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "144289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8521,
            "unit": "ns/op\t    5987 B/op\t      75 allocs/op",
            "extra": "140886 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8521,
            "unit": "ns/op",
            "extra": "140886 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5987,
            "unit": "B/op",
            "extra": "140886 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "140886 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8250,
            "unit": "ns/op\t    5991 B/op\t      76 allocs/op",
            "extra": "142712 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8250,
            "unit": "ns/op",
            "extra": "142712 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5991,
            "unit": "B/op",
            "extra": "142712 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "142712 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 8009,
            "unit": "ns/op\t    5990 B/op\t      76 allocs/op",
            "extra": "149325 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 8009,
            "unit": "ns/op",
            "extra": "149325 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5990,
            "unit": "B/op",
            "extra": "149325 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "149325 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8286,
            "unit": "ns/op\t    5987 B/op\t      76 allocs/op",
            "extra": "144212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8286,
            "unit": "ns/op",
            "extra": "144212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5987,
            "unit": "B/op",
            "extra": "144212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "144212 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9843,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "122728 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9843,
            "unit": "ns/op",
            "extra": "122728 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "122728 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "122728 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 19170,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "62190 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 19170,
            "unit": "ns/op",
            "extra": "62190 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "62190 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "62190 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2107,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "498039 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2107,
            "unit": "ns/op",
            "extra": "498039 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "498039 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "498039 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8765,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "138804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8765,
            "unit": "ns/op",
            "extra": "138804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "138804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "138804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5450,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "215934 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5450,
            "unit": "ns/op",
            "extra": "215934 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "215934 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "215934 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 43127,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 43127,
            "unit": "ns/op",
            "extra": "28105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28105 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39589,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30288 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39589,
            "unit": "ns/op",
            "extra": "30288 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30288 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30288 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3844,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "312646 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3844,
            "unit": "ns/op",
            "extra": "312646 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "312646 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "312646 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1829,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "631148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1829,
            "unit": "ns/op",
            "extra": "631148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "631148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "631148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1850,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "630902 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1850,
            "unit": "ns/op",
            "extra": "630902 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "630902 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "630902 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2225,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "516457 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2225,
            "unit": "ns/op",
            "extra": "516457 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "516457 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "516457 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2277,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "523494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2277,
            "unit": "ns/op",
            "extra": "523494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "523494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "523494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2179,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "528552 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2179,
            "unit": "ns/op",
            "extra": "528552 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "528552 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "528552 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3204,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "364134 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3204,
            "unit": "ns/op",
            "extra": "364134 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "364134 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "364134 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6126,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "196989 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6126,
            "unit": "ns/op",
            "extra": "196989 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "196989 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "196989 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6232,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "195774 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6232,
            "unit": "ns/op",
            "extra": "195774 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "195774 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "195774 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6213,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "190480 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6213,
            "unit": "ns/op",
            "extra": "190480 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "190480 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "190480 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12253,
            "unit": "ns/op\t    8327 B/op\t     117 allocs/op",
            "extra": "99307 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12253,
            "unit": "ns/op",
            "extra": "99307 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8327,
            "unit": "B/op",
            "extra": "99307 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "99307 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6920,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "169916 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6920,
            "unit": "ns/op",
            "extra": "169916 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "169916 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "169916 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22935,
            "unit": "ns/op\t   13937 B/op\t     217 allocs/op",
            "extra": "52878 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22935,
            "unit": "ns/op",
            "extra": "52878 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13937,
            "unit": "B/op",
            "extra": "52878 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "52878 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2832,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "415150 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2832,
            "unit": "ns/op",
            "extra": "415150 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "415150 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "415150 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2848,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "419334 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2848,
            "unit": "ns/op",
            "extra": "419334 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "419334 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "419334 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2842,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "416742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2842,
            "unit": "ns/op",
            "extra": "416742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "416742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "416742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2827,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "424417 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2827,
            "unit": "ns/op",
            "extra": "424417 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "424417 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "424417 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2834,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "414873 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2834,
            "unit": "ns/op",
            "extra": "414873 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "414873 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "414873 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 635480,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1869 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 635480,
            "unit": "ns/op",
            "extra": "1869 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1869 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1869 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 516630,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 516630,
            "unit": "ns/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 607.7,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1971612 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 607.7,
            "unit": "ns/op",
            "extra": "1971612 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1971612 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1971612 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 608.1,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1948316 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 608.1,
            "unit": "ns/op",
            "extra": "1948316 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1948316 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1948316 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 619.4,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1943306 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 619.4,
            "unit": "ns/op",
            "extra": "1943306 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1943306 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1943306 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 160.7,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8412930 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 160.7,
            "unit": "ns/op",
            "extra": "8412930 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8412930 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8412930 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 7929,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "150373 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 7929,
            "unit": "ns/op",
            "extra": "150373 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "150373 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "150373 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "richard.wooding@spandigital.com",
            "name": "Richard Wooding",
            "username": "richardwooding"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4e367beed4292258c4d7e31e8f4a32d111996ac3",
          "message": "fix: Add byte array length validation to prevent resource exhaustion (fixes #36) (#82)\n\n## Summary\nAdds validation to prevent large byte arrays from causing memory exhaustion\nthrough hex-encoded SQL string expansion in non-parameterized mode.\n\n## Changes\n- Added `maxByteArrayLength = 10,000` constant (≈40KB SQL output)\n- Implemented length validation in `visitConst()` BytesValue case\n- Returns descriptive error when limit exceeded\n- Parameterized mode bypasses this check (bytes sent directly to DB driver)\n\n## Security Context\n- Addresses CWE-400 (Uncontrolled Resource Consumption)\n- Each byte expands to ~4 characters in hex format (\\xDE)\n- 10,000 bytes → ~40KB SQL (acceptable)\n- 100,000 bytes → ~400KB SQL (now prevented)\n\n## Testing\n- Added comprehensive test suite in byte_array_test.go\n- Tests cover small arrays, parameterized mode, hex encoding, schemas\n- All existing tests pass\n- Includes documentation tests explaining the feature\n\n## Documentation Updates\n- README.md: Added byte array limits to security features\n- CLAUDE.md: Updated resource exhaustion protection section\n\n## Alignment with Existing Security Features\nFollows same pattern as:\n- maxComprehensionDepth = 3 (issue #35)\n- defaultMaxRecursionDepth = 100 (issue #27)\n- defaultMaxSQLOutputLength = 50,000 (issue #33)\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-authored-by: Claude <noreply@anthropic.com>",
          "timestamp": "2025-10-30T23:55:08+02:00",
          "tree_id": "9406f5b22213bc1bf6927b13f8cfd2381c41ceb7",
          "url": "https://github.com/SPANDigital/cel2sql/commit/4e367beed4292258c4d7e31e8f4a32d111996ac3"
        },
        "date": 1761861443021,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkByteArrayConversion",
            "value": 8019,
            "unit": "ns/op\t   14708 B/op\t      30 allocs/op",
            "extra": "152330 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - ns/op",
            "value": 8019,
            "unit": "ns/op",
            "extra": "152330 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - B/op",
            "value": 14708,
            "unit": "B/op",
            "extra": "152330 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "152330 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2009,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "571363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2009,
            "unit": "ns/op",
            "extra": "571363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "571363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "571363 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2025,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "588152 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2025,
            "unit": "ns/op",
            "extra": "588152 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "588152 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "588152 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2427,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "541556 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2427,
            "unit": "ns/op",
            "extra": "541556 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "541556 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "541556 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1381,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "870955 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1381,
            "unit": "ns/op",
            "extra": "870955 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "870955 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "870955 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 3007,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "380074 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 3007,
            "unit": "ns/op",
            "extra": "380074 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "380074 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "380074 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3648,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "326289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3648,
            "unit": "ns/op",
            "extra": "326289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "326289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "326289 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2734,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "423406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2734,
            "unit": "ns/op",
            "extra": "423406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "423406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "423406 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2782,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "440305 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2782,
            "unit": "ns/op",
            "extra": "440305 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "440305 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "440305 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3141,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "372864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3141,
            "unit": "ns/op",
            "extra": "372864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "372864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "372864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3165,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "378772 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3165,
            "unit": "ns/op",
            "extra": "378772 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "378772 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "378772 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2698,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "442268 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2698,
            "unit": "ns/op",
            "extra": "442268 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "442268 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "442268 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3075,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "389484 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3075,
            "unit": "ns/op",
            "extra": "389484 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "389484 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "389484 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7374,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "160384 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7374,
            "unit": "ns/op",
            "extra": "160384 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "160384 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "160384 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5247,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "227640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5247,
            "unit": "ns/op",
            "extra": "227640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "227640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "227640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7617,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "155764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7617,
            "unit": "ns/op",
            "extra": "155764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "155764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "155764 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5852,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "202208 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5852,
            "unit": "ns/op",
            "extra": "202208 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "202208 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "202208 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7576,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "159226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7576,
            "unit": "ns/op",
            "extra": "159226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "159226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "159226 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 7007,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "172339 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 7007,
            "unit": "ns/op",
            "extra": "172339 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "172339 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "172339 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7268,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "163509 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7268,
            "unit": "ns/op",
            "extra": "163509 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "163509 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "163509 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 6053,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "204483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 6053,
            "unit": "ns/op",
            "extra": "204483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "204483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "204483 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3505,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "330318 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3505,
            "unit": "ns/op",
            "extra": "330318 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "330318 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "330318 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4107,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "286640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4107,
            "unit": "ns/op",
            "extra": "286640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "286640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "286640 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2788,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "430846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2788,
            "unit": "ns/op",
            "extra": "430846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "430846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "430846 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3249,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "370004 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3249,
            "unit": "ns/op",
            "extra": "370004 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "370004 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "370004 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3500,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "339122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3500,
            "unit": "ns/op",
            "extra": "339122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "339122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "339122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7490,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "157972 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7490,
            "unit": "ns/op",
            "extra": "157972 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "157972 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "157972 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8275,
            "unit": "ns/op\t    5970 B/op\t      75 allocs/op",
            "extra": "143758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8275,
            "unit": "ns/op",
            "extra": "143758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5970,
            "unit": "B/op",
            "extra": "143758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "143758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7877,
            "unit": "ns/op\t    5972 B/op\t      75 allocs/op",
            "extra": "154500 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7877,
            "unit": "ns/op",
            "extra": "154500 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5972,
            "unit": "B/op",
            "extra": "154500 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "154500 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8334,
            "unit": "ns/op\t    5995 B/op\t      75 allocs/op",
            "extra": "142338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8334,
            "unit": "ns/op",
            "extra": "142338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5995,
            "unit": "B/op",
            "extra": "142338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8088,
            "unit": "ns/op\t    5990 B/op\t      76 allocs/op",
            "extra": "149104 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8088,
            "unit": "ns/op",
            "extra": "149104 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5990,
            "unit": "B/op",
            "extra": "149104 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "149104 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7906,
            "unit": "ns/op\t    5997 B/op\t      76 allocs/op",
            "extra": "152599 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7906,
            "unit": "ns/op",
            "extra": "152599 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5997,
            "unit": "B/op",
            "extra": "152599 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "152599 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8164,
            "unit": "ns/op\t    5982 B/op\t      76 allocs/op",
            "extra": "145597 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8164,
            "unit": "ns/op",
            "extra": "145597 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5982,
            "unit": "B/op",
            "extra": "145597 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "145597 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9781,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "121726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9781,
            "unit": "ns/op",
            "extra": "121726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "121726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "121726 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 18931,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "63007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 18931,
            "unit": "ns/op",
            "extra": "63007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "63007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "63007 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2102,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "566862 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2102,
            "unit": "ns/op",
            "extra": "566862 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "566862 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "566862 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8764,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "137490 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8764,
            "unit": "ns/op",
            "extra": "137490 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "137490 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "137490 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5457,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "218187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5457,
            "unit": "ns/op",
            "extra": "218187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "218187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "218187 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42287,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28111 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42287,
            "unit": "ns/op",
            "extra": "28111 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28111 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28111 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39356,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39356,
            "unit": "ns/op",
            "extra": "30604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30604 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3766,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "314656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3766,
            "unit": "ns/op",
            "extra": "314656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "314656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "314656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1818,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "629295 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1818,
            "unit": "ns/op",
            "extra": "629295 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "629295 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "629295 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1851,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "642840 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1851,
            "unit": "ns/op",
            "extra": "642840 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "642840 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "642840 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2210,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "526370 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2210,
            "unit": "ns/op",
            "extra": "526370 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "526370 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "526370 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2259,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "528450 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2259,
            "unit": "ns/op",
            "extra": "528450 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "528450 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "528450 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2169,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "549793 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2169,
            "unit": "ns/op",
            "extra": "549793 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "549793 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "549793 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3195,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "370686 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3195,
            "unit": "ns/op",
            "extra": "370686 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "370686 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "370686 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6101,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "195315 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6101,
            "unit": "ns/op",
            "extra": "195315 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "195315 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "195315 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6193,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "193825 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6193,
            "unit": "ns/op",
            "extra": "193825 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "193825 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193825 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6159,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "194330 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6159,
            "unit": "ns/op",
            "extra": "194330 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "194330 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "194330 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12132,
            "unit": "ns/op\t    8315 B/op\t     117 allocs/op",
            "extra": "93043 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12132,
            "unit": "ns/op",
            "extra": "93043 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8315,
            "unit": "B/op",
            "extra": "93043 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "93043 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6889,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "175034 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6889,
            "unit": "ns/op",
            "extra": "175034 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "175034 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "175034 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22527,
            "unit": "ns/op\t   13949 B/op\t     217 allocs/op",
            "extra": "53468 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22527,
            "unit": "ns/op",
            "extra": "53468 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13949,
            "unit": "B/op",
            "extra": "53468 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "53468 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2813,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "414156 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2813,
            "unit": "ns/op",
            "extra": "414156 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "414156 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "414156 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2817,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "429751 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2817,
            "unit": "ns/op",
            "extra": "429751 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "429751 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "429751 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2861,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "420841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2861,
            "unit": "ns/op",
            "extra": "420841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "420841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "420841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2800,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "413660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2800,
            "unit": "ns/op",
            "extra": "413660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "413660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "413660 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2816,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "419931 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2816,
            "unit": "ns/op",
            "extra": "419931 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "419931 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "419931 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 614551,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1933 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 614551,
            "unit": "ns/op",
            "extra": "1933 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1933 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1933 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 509793,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2372 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 509793,
            "unit": "ns/op",
            "extra": "2372 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2372 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2372 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 603.5,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1977230 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 603.5,
            "unit": "ns/op",
            "extra": "1977230 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1977230 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1977230 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 608.6,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1978456 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 608.6,
            "unit": "ns/op",
            "extra": "1978456 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1978456 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1978456 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 611.6,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1950904 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 611.6,
            "unit": "ns/op",
            "extra": "1950904 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1950904 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1950904 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 143.6,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8397595 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 143.6,
            "unit": "ns/op",
            "extra": "8397595 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8397595 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8397595 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 7877,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "151388 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 7877,
            "unit": "ns/op",
            "extra": "151388 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "151388 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "151388 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "SPANDigital",
            "username": "SPANDigital"
          },
          "committer": {
            "name": "SPANDigital",
            "username": "SPANDigital"
          },
          "id": "223ada38073a83bdb738a6c32bf9aa5293ce5e8e",
          "message": "fix: Standardize error message format throughout codebase (fixes #38)",
          "timestamp": "2025-10-30T21:55:12Z",
          "url": "https://github.com/SPANDigital/cel2sql/pull/83/commits/223ada38073a83bdb738a6c32bf9aa5293ce5e8e"
        },
        "date": 1761862996776,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkByteArrayConversion",
            "value": 7952,
            "unit": "ns/op\t   14708 B/op\t      30 allocs/op",
            "extra": "150984 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - ns/op",
            "value": 7952,
            "unit": "ns/op",
            "extra": "150984 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - B/op",
            "value": 14708,
            "unit": "B/op",
            "extra": "150984 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "150984 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2003,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "572251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2003,
            "unit": "ns/op",
            "extra": "572251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "572251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "572251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2023,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "590842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2023,
            "unit": "ns/op",
            "extra": "590842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "590842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "590842 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2150,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "530608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2150,
            "unit": "ns/op",
            "extra": "530608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "530608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "530608 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1507,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "853681 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1507,
            "unit": "ns/op",
            "extra": "853681 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "853681 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "853681 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 2998,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "385351 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 2998,
            "unit": "ns/op",
            "extra": "385351 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "385351 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "385351 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3672,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "323462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3672,
            "unit": "ns/op",
            "extra": "323462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "323462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "323462 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2749,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "373942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2749,
            "unit": "ns/op",
            "extra": "373942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "373942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "373942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2743,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "432550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2743,
            "unit": "ns/op",
            "extra": "432550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "432550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "432550 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3100,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "384794 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3100,
            "unit": "ns/op",
            "extra": "384794 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "384794 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "384794 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3105,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "364023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3105,
            "unit": "ns/op",
            "extra": "364023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "364023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "364023 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2709,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "440409 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2709,
            "unit": "ns/op",
            "extra": "440409 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "440409 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "440409 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3074,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "394804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3074,
            "unit": "ns/op",
            "extra": "394804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "394804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "394804 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7309,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "160496 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7309,
            "unit": "ns/op",
            "extra": "160496 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "160496 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "160496 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5187,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "227792 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5187,
            "unit": "ns/op",
            "extra": "227792 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "227792 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "227792 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7587,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7587,
            "unit": "ns/op",
            "extra": "158742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158742 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5825,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "202790 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5825,
            "unit": "ns/op",
            "extra": "202790 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "202790 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "202790 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7593,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "159519 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7593,
            "unit": "ns/op",
            "extra": "159519 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "159519 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "159519 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6951,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "169182 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6951,
            "unit": "ns/op",
            "extra": "169182 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "169182 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "169182 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7330,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "164161 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7330,
            "unit": "ns/op",
            "extra": "164161 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "164161 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "164161 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 5884,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "170006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 5884,
            "unit": "ns/op",
            "extra": "170006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "170006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "170006 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3484,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "335949 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3484,
            "unit": "ns/op",
            "extra": "335949 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "335949 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "335949 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4081,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "294946 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4081,
            "unit": "ns/op",
            "extra": "294946 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "294946 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "294946 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2831,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "420910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2831,
            "unit": "ns/op",
            "extra": "420910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "420910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "420910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3233,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "368512 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3233,
            "unit": "ns/op",
            "extra": "368512 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "368512 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "368512 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3461,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "339891 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3461,
            "unit": "ns/op",
            "extra": "339891 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "339891 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "339891 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7600,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "158960 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7600,
            "unit": "ns/op",
            "extra": "158960 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "158960 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "158960 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8363,
            "unit": "ns/op\t    5979 B/op\t      75 allocs/op",
            "extra": "144453 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8363,
            "unit": "ns/op",
            "extra": "144453 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5979,
            "unit": "B/op",
            "extra": "144453 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "144453 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7962,
            "unit": "ns/op\t    5970 B/op\t      75 allocs/op",
            "extra": "151516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7962,
            "unit": "ns/op",
            "extra": "151516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5970,
            "unit": "B/op",
            "extra": "151516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "151516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8388,
            "unit": "ns/op\t    5995 B/op\t      75 allocs/op",
            "extra": "142976 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8388,
            "unit": "ns/op",
            "extra": "142976 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5995,
            "unit": "B/op",
            "extra": "142976 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142976 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8351,
            "unit": "ns/op\t    5994 B/op\t      76 allocs/op",
            "extra": "149461 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8351,
            "unit": "ns/op",
            "extra": "149461 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5994,
            "unit": "B/op",
            "extra": "149461 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "149461 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 8052,
            "unit": "ns/op\t    5994 B/op\t      76 allocs/op",
            "extra": "147075 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 8052,
            "unit": "ns/op",
            "extra": "147075 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5994,
            "unit": "B/op",
            "extra": "147075 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147075 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8225,
            "unit": "ns/op\t    5980 B/op\t      76 allocs/op",
            "extra": "147338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8225,
            "unit": "ns/op",
            "extra": "147338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5980,
            "unit": "B/op",
            "extra": "147338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147338 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9809,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "122095 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9809,
            "unit": "ns/op",
            "extra": "122095 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "122095 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "122095 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 18930,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "62438 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 18930,
            "unit": "ns/op",
            "extra": "62438 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "62438 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "62438 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2086,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "549531 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2086,
            "unit": "ns/op",
            "extra": "549531 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "549531 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "549531 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8824,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "138242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8824,
            "unit": "ns/op",
            "extra": "138242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "138242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "138242 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5481,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "217155 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5481,
            "unit": "ns/op",
            "extra": "217155 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "217155 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "217155 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 43976,
            "unit": "ns/op\t   29107 B/op\t     424 allocs/op",
            "extra": "27823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 43976,
            "unit": "ns/op",
            "extra": "27823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29107,
            "unit": "B/op",
            "extra": "27823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "27823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 40640,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "29575 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 40640,
            "unit": "ns/op",
            "extra": "29575 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "29575 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "29575 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3770,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "314557 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3770,
            "unit": "ns/op",
            "extra": "314557 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "314557 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "314557 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1858,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "626863 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1858,
            "unit": "ns/op",
            "extra": "626863 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "626863 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "626863 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1878,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "638634 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1878,
            "unit": "ns/op",
            "extra": "638634 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "638634 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "638634 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2276,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "515942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2276,
            "unit": "ns/op",
            "extra": "515942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "515942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "515942 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2294,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "511594 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2294,
            "unit": "ns/op",
            "extra": "511594 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "511594 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "511594 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2201,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "555470 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2201,
            "unit": "ns/op",
            "extra": "555470 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "555470 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "555470 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3188,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "350415 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3188,
            "unit": "ns/op",
            "extra": "350415 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "350415 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "350415 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6102,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "194703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6102,
            "unit": "ns/op",
            "extra": "194703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "194703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "194703 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6136,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "195231 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6136,
            "unit": "ns/op",
            "extra": "195231 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "195231 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "195231 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6135,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "196005 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6135,
            "unit": "ns/op",
            "extra": "196005 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "196005 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "196005 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12121,
            "unit": "ns/op\t    8336 B/op\t     117 allocs/op",
            "extra": "99960 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12121,
            "unit": "ns/op",
            "extra": "99960 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8336,
            "unit": "B/op",
            "extra": "99960 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "99960 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6876,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "174009 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6876,
            "unit": "ns/op",
            "extra": "174009 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "174009 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "174009 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22706,
            "unit": "ns/op\t   13938 B/op\t     217 allocs/op",
            "extra": "51361 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22706,
            "unit": "ns/op",
            "extra": "51361 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13938,
            "unit": "B/op",
            "extra": "51361 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "51361 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2758,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "413919 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2758,
            "unit": "ns/op",
            "extra": "413919 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "413919 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "413919 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2822,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "430183 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2822,
            "unit": "ns/op",
            "extra": "430183 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "430183 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "430183 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2782,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "428864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2782,
            "unit": "ns/op",
            "extra": "428864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "428864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "428864 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2792,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "428649 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2792,
            "unit": "ns/op",
            "extra": "428649 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "428649 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "428649 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2805,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "412630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2805,
            "unit": "ns/op",
            "extra": "412630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "412630 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "412630 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 618417,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1917 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 618417,
            "unit": "ns/op",
            "extra": "1917 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1917 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1917 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 517826,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 517826,
            "unit": "ns/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 605.7,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1968246 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 605.7,
            "unit": "ns/op",
            "extra": "1968246 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1968246 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1968246 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 658.9,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1966528 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 658.9,
            "unit": "ns/op",
            "extra": "1966528 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1966528 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1966528 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 607.5,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1971654 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 607.5,
            "unit": "ns/op",
            "extra": "1971654 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1971654 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1971654 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 143.9,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8279176 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 143.9,
            "unit": "ns/op",
            "extra": "8279176 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8279176 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8279176 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 7940,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "152191 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 7940,
            "unit": "ns/op",
            "extra": "152191 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "152191 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "152191 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "richard.wooding@spandigital.com",
            "name": "Richard Wooding",
            "username": "richardwooding"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "65f8003eacd71bca6bd0b1ebe1391efed4b48952",
          "message": "fix: Standardize error message format throughout codebase (fixes #38) (#83)\n\nThis commit addresses issue #38 by standardizing error handling across\nthe cel2sql codebase with the following improvements:\n\n## Changes\n\n### Added Sentinel Errors\n- Introduced exported sentinel errors in errors.go to enable `errors.Is()` checking:\n  - `ErrUnsupportedExpression`, `ErrInvalidFieldName`, `ErrInvalidSchema`\n  - `ErrInvalidRegexPattern`, `ErrMaxDepthExceeded`, `ErrMaxOutputLengthExceeded`\n  - `ErrInvalidComprehension`, `ErrMaxComprehensionDepthExceeded`\n  - `ErrInvalidArguments`, `ErrInvalidTimestampOperation`, `ErrInvalidDuration`\n  - `ErrInvalidJSONPath`, `ErrInvalidOperator`, `ErrUnsupportedType`\n  - `ErrContextCanceled`, `ErrInvalidByteArrayLength`\n- Added `pg.ErrInvalidSchema` for PostgreSQL provider package\n\n### Improved Error Wrapping\n- All errors now use `fmt.Errorf()` with `%w` for proper error wrapping\n- Sentinel errors provide context while maintaining error chains\n- Enables better error handling with `errors.Is()` and `errors.As()`\n\n### Enhanced Error Messages\n- Added operation context to error messages (e.g., \"failed to identify comprehension\")\n- Improved specificity for debugging (e.g., \"depth 4 exceeds limit of 3\")\n- Maintained security-conscious error handling in pg/provider.go (no credential exposure)\n\n### Updated Tests\n- Modified tests to check for new error message formats\n- Changed from exact string matching to `Contains()` for flexibility\n- Updated error assertions to work with wrapped sentinel errors\n\n## Benefits\n- Better debugging with improved error context\n- Enables programmatic error handling with `errors.Is()`\n- More maintainable error handling patterns\n- Consistent error format across the codebase\n\n🤖 Generated with [Claude Code](https://claude.com/claude-code)\n\nCo-authored-by: Claude <noreply@anthropic.com>",
          "timestamp": "2025-10-31T00:27:14+02:00",
          "tree_id": "ea3f301b15739f074fee1ada334a1a5bba4e6d06",
          "url": "https://github.com/SPANDigital/cel2sql/commit/65f8003eacd71bca6bd0b1ebe1391efed4b48952"
        },
        "date": 1761863382253,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkByteArrayConversion",
            "value": 8130,
            "unit": "ns/op\t   14708 B/op\t      30 allocs/op",
            "extra": "150963 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - ns/op",
            "value": 8130,
            "unit": "ns/op",
            "extra": "150963 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - B/op",
            "value": 14708,
            "unit": "B/op",
            "extra": "150963 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "150963 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2004,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "567714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2004,
            "unit": "ns/op",
            "extra": "567714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "567714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "567714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2022,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "589675 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2022,
            "unit": "ns/op",
            "extra": "589675 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "589675 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "589675 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2165,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "539703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2165,
            "unit": "ns/op",
            "extra": "539703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "539703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "539703 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1424,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "833744 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1424,
            "unit": "ns/op",
            "extra": "833744 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "833744 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "833744 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 3112,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "390148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 3112,
            "unit": "ns/op",
            "extra": "390148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "390148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "390148 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 4145,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "252738 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 4145,
            "unit": "ns/op",
            "extra": "252738 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "252738 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "252738 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2746,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "411369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2746,
            "unit": "ns/op",
            "extra": "411369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "411369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "411369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2733,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "429013 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2733,
            "unit": "ns/op",
            "extra": "429013 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "429013 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "429013 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul",
            "value": 3120,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "379516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3120,
            "unit": "ns/op",
            "extra": "379516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "379516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "379516 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3124,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "381570 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3124,
            "unit": "ns/op",
            "extra": "381570 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "381570 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "381570 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2741,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "429036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2741,
            "unit": "ns/op",
            "extra": "429036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "429036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "429036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3067,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "393020 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3067,
            "unit": "ns/op",
            "extra": "393020 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "393020 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "393020 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7609,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "158313 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7609,
            "unit": "ns/op",
            "extra": "158313 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "158313 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "158313 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5217,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "226251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5217,
            "unit": "ns/op",
            "extra": "226251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "226251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "226251 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7520,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158833 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7520,
            "unit": "ns/op",
            "extra": "158833 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158833 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158833 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5851,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "205330 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5851,
            "unit": "ns/op",
            "extra": "205330 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "205330 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "205330 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7628,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "157497 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7628,
            "unit": "ns/op",
            "extra": "157497 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "157497 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "157497 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6916,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "171829 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6916,
            "unit": "ns/op",
            "extra": "171829 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "171829 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "171829 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7353,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "162831 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7353,
            "unit": "ns/op",
            "extra": "162831 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "162831 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "162831 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 5896,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "202028 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 5896,
            "unit": "ns/op",
            "extra": "202028 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "202028 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "202028 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3476,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "338138 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3476,
            "unit": "ns/op",
            "extra": "338138 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "338138 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "338138 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4101,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "289998 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4101,
            "unit": "ns/op",
            "extra": "289998 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "289998 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "289998 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2774,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "419328 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2774,
            "unit": "ns/op",
            "extra": "419328 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "419328 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "419328 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3229,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "370791 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3229,
            "unit": "ns/op",
            "extra": "370791 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "370791 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "370791 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3453,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "343918 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3453,
            "unit": "ns/op",
            "extra": "343918 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "343918 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "343918 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7479,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "158610 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7479,
            "unit": "ns/op",
            "extra": "158610 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "158610 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "158610 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8290,
            "unit": "ns/op\t    5975 B/op\t      75 allocs/op",
            "extra": "144823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8290,
            "unit": "ns/op",
            "extra": "144823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5975,
            "unit": "B/op",
            "extra": "144823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "144823 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7942,
            "unit": "ns/op\t    5976 B/op\t      75 allocs/op",
            "extra": "149689 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7942,
            "unit": "ns/op",
            "extra": "149689 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5976,
            "unit": "B/op",
            "extra": "149689 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "149689 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8444,
            "unit": "ns/op\t    5989 B/op\t      75 allocs/op",
            "extra": "142518 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8444,
            "unit": "ns/op",
            "extra": "142518 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5989,
            "unit": "B/op",
            "extra": "142518 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142518 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8088,
            "unit": "ns/op\t    5983 B/op\t      76 allocs/op",
            "extra": "148081 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8088,
            "unit": "ns/op",
            "extra": "148081 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5983,
            "unit": "B/op",
            "extra": "148081 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "148081 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7961,
            "unit": "ns/op\t    5994 B/op\t      76 allocs/op",
            "extra": "147373 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7961,
            "unit": "ns/op",
            "extra": "147373 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5994,
            "unit": "B/op",
            "extra": "147373 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147373 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8201,
            "unit": "ns/op\t    5984 B/op\t      76 allocs/op",
            "extra": "145788 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8201,
            "unit": "ns/op",
            "extra": "145788 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5984,
            "unit": "B/op",
            "extra": "145788 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "145788 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9816,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "121563 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9816,
            "unit": "ns/op",
            "extra": "121563 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "121563 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "121563 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 18853,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "63103 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 18853,
            "unit": "ns/op",
            "extra": "63103 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "63103 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "63103 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2094,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "556494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2094,
            "unit": "ns/op",
            "extra": "556494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "556494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "556494 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8677,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "137904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8677,
            "unit": "ns/op",
            "extra": "137904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "137904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "137904 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5442,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "219678 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5442,
            "unit": "ns/op",
            "extra": "219678 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "219678 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "219678 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42689,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "28395 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42689,
            "unit": "ns/op",
            "extra": "28395 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "28395 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "28395 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39490,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30368 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39490,
            "unit": "ns/op",
            "extra": "30368 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30368 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30368 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3752,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "308120 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3752,
            "unit": "ns/op",
            "extra": "308120 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "308120 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "308120 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1806,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "661585 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1806,
            "unit": "ns/op",
            "extra": "661585 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "661585 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "661585 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1837,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "648122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1837,
            "unit": "ns/op",
            "extra": "648122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "648122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "648122 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2215,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "532951 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2215,
            "unit": "ns/op",
            "extra": "532951 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "532951 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "532951 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2256,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "524850 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2256,
            "unit": "ns/op",
            "extra": "524850 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "524850 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "524850 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2159,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "552196 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2159,
            "unit": "ns/op",
            "extra": "552196 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "552196 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "552196 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3181,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "377996 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3181,
            "unit": "ns/op",
            "extra": "377996 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "377996 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "377996 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6154,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "193502 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6154,
            "unit": "ns/op",
            "extra": "193502 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "193502 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193502 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6138,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "194698 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6138,
            "unit": "ns/op",
            "extra": "194698 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "194698 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "194698 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6104,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "196429 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6104,
            "unit": "ns/op",
            "extra": "196429 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "196429 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "196429 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12030,
            "unit": "ns/op\t    8323 B/op\t     117 allocs/op",
            "extra": "96818 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12030,
            "unit": "ns/op",
            "extra": "96818 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8323,
            "unit": "B/op",
            "extra": "96818 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "96818 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6832,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "176445 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6832,
            "unit": "ns/op",
            "extra": "176445 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "176445 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "176445 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22405,
            "unit": "ns/op\t   13930 B/op\t     217 allocs/op",
            "extra": "52821 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22405,
            "unit": "ns/op",
            "extra": "52821 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13930,
            "unit": "B/op",
            "extra": "52821 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "52821 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2778,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "417099 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2778,
            "unit": "ns/op",
            "extra": "417099 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "417099 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "417099 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2788,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "424435 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2788,
            "unit": "ns/op",
            "extra": "424435 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "424435 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "424435 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2776,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "423753 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2776,
            "unit": "ns/op",
            "extra": "423753 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "423753 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "423753 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2778,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "433568 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2778,
            "unit": "ns/op",
            "extra": "433568 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "433568 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "433568 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2791,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "420880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2791,
            "unit": "ns/op",
            "extra": "420880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "420880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "420880 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 619112,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1909 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 619112,
            "unit": "ns/op",
            "extra": "1909 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1909 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1909 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 516226,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2317 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 516226,
            "unit": "ns/op",
            "extra": "2317 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2317 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2317 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 602.9,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1961288 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 602.9,
            "unit": "ns/op",
            "extra": "1961288 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1961288 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1961288 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 604,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1979431 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 604,
            "unit": "ns/op",
            "extra": "1979431 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1979431 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1979431 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 659.2,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1802524 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 659.2,
            "unit": "ns/op",
            "extra": "1802524 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1802524 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1802524 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 148.6,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8167635 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 148.6,
            "unit": "ns/op",
            "extra": "8167635 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8167635 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8167635 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 7914,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "150018 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 7914,
            "unit": "ns/op",
            "extra": "150018 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "150018 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "150018 times\n4 procs"
          }
        ]
      },
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
          "id": "8b6387d7e8877844831fedada363cc4fbd0d34ba",
          "message": "chore: Update CHANGELOG for v3.2.0 release",
          "timestamp": "2025-10-31T00:29:21+02:00",
          "tree_id": "8433b74db593c36cd5a2d640e173e6a377521701",
          "url": "https://github.com/SPANDigital/cel2sql/commit/8b6387d7e8877844831fedada363cc4fbd0d34ba"
        },
        "date": 1761863510787,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkByteArrayConversion",
            "value": 9425,
            "unit": "ns/op\t   14708 B/op\t      30 allocs/op",
            "extra": "151172 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - ns/op",
            "value": 9425,
            "unit": "ns/op",
            "extra": "151172 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - B/op",
            "value": 14708,
            "unit": "B/op",
            "extra": "151172 times\n4 procs"
          },
          {
            "name": "BenchmarkByteArrayConversion - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "151172 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality",
            "value": 2025,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "573386 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - ns/op",
            "value": 2025,
            "unit": "ns/op",
            "extra": "573386 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "573386 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/equality - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "573386 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than",
            "value": 2046,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "564050 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - ns/op",
            "value": 2046,
            "unit": "ns/op",
            "extra": "564050 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "564050 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/greater_than - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "564050 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality",
            "value": 2212,
            "unit": "ns/op\t    1632 B/op\t      27 allocs/op",
            "extra": "534624 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - ns/op",
            "value": 2212,
            "unit": "ns/op",
            "extra": "534624 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - B/op",
            "value": 1632,
            "unit": "B/op",
            "extra": "534624 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/string_equality - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "534624 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check",
            "value": 1379,
            "unit": "ns/op\t    1168 B/op\t      16 allocs/op",
            "extra": "810436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - ns/op",
            "value": 1379,
            "unit": "ns/op",
            "extra": "810436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - B/op",
            "value": 1168,
            "unit": "B/op",
            "extra": "810436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertSimple/boolean_check - allocs/op",
            "value": 16,
            "unit": "allocs/op",
            "extra": "810436 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and",
            "value": 3019,
            "unit": "ns/op\t    2152 B/op\t      36 allocs/op",
            "extra": "383892 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - ns/op",
            "value": 3019,
            "unit": "ns/op",
            "extra": "383892 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - B/op",
            "value": 2152,
            "unit": "B/op",
            "extra": "383892 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_and - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "383892 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or",
            "value": 3675,
            "unit": "ns/op\t    2576 B/op\t      45 allocs/op",
            "extra": "318304 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - ns/op",
            "value": 3675,
            "unit": "ns/op",
            "extra": "318304 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - B/op",
            "value": 2576,
            "unit": "B/op",
            "extra": "318304 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/logical_or - allocs/op",
            "value": 45,
            "unit": "allocs/op",
            "extra": "318304 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add",
            "value": 2796,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "426308 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - ns/op",
            "value": 2796,
            "unit": "ns/op",
            "extra": "426308 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "426308 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_add - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "426308 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub",
            "value": 2775,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "431342 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_sub - ns/op",
            "value": 2775,
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
            "value": 3206,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "371506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - ns/op",
            "value": 3206,
            "unit": "ns/op",
            "extra": "371506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "371506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_mul - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "371506 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div",
            "value": 3148,
            "unit": "ns/op\t    2128 B/op\t      40 allocs/op",
            "extra": "376309 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - ns/op",
            "value": 3148,
            "unit": "ns/op",
            "extra": "376309 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - B/op",
            "value": 2128,
            "unit": "B/op",
            "extra": "376309 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/arithmetic_div - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "376309 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo",
            "value": 2721,
            "unit": "ns/op\t    2032 B/op\t      35 allocs/op",
            "extra": "441910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - ns/op",
            "value": 2721,
            "unit": "ns/op",
            "extra": "441910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - B/op",
            "value": 2032,
            "unit": "B/op",
            "extra": "441910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/modulo - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "441910 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat",
            "value": 3070,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "391903 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - ns/op",
            "value": 3070,
            "unit": "ns/op",
            "extra": "391903 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "391903 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/string_concat - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "391903 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression",
            "value": 7357,
            "unit": "ns/op\t    4689 B/op\t      82 allocs/op",
            "extra": "161697 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - ns/op",
            "value": 7357,
            "unit": "ns/op",
            "extra": "161697 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - B/op",
            "value": 4689,
            "unit": "B/op",
            "extra": "161697 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertOperators/complex_expression - allocs/op",
            "value": 82,
            "unit": "allocs/op",
            "extra": "161697 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple",
            "value": 5234,
            "unit": "ns/op\t    4089 B/op\t      67 allocs/op",
            "extra": "229297 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - ns/op",
            "value": 5234,
            "unit": "ns/op",
            "extra": "229297 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - B/op",
            "value": 4089,
            "unit": "B/op",
            "extra": "229297 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_simple - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "229297 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex",
            "value": 7538,
            "unit": "ns/op\t    5914 B/op\t      89 allocs/op",
            "extra": "158758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - ns/op",
            "value": 7538,
            "unit": "ns/op",
            "extra": "158758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - B/op",
            "value": 5914,
            "unit": "B/op",
            "extra": "158758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/all_complex - allocs/op",
            "value": 89,
            "unit": "allocs/op",
            "extra": "158758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple",
            "value": 5990,
            "unit": "ns/op\t    4634 B/op\t      76 allocs/op",
            "extra": "202896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - ns/op",
            "value": 5990,
            "unit": "ns/op",
            "extra": "202896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - B/op",
            "value": 4634,
            "unit": "B/op",
            "extra": "202896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_simple - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "202896 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex",
            "value": 7676,
            "unit": "ns/op\t    6155 B/op\t      95 allocs/op",
            "extra": "156399 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - ns/op",
            "value": 7676,
            "unit": "ns/op",
            "extra": "156399 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - B/op",
            "value": 6155,
            "unit": "B/op",
            "extra": "156399 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_complex - allocs/op",
            "value": 95,
            "unit": "allocs/op",
            "extra": "156399 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one",
            "value": 6959,
            "unit": "ns/op\t    5673 B/op\t      88 allocs/op",
            "extra": "168162 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - ns/op",
            "value": 6959,
            "unit": "ns/op",
            "extra": "168162 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - B/op",
            "value": 5673,
            "unit": "B/op",
            "extra": "168162 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/exists_one - allocs/op",
            "value": 88,
            "unit": "allocs/op",
            "extra": "168162 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter",
            "value": 7281,
            "unit": "ns/op\t    5842 B/op\t     105 allocs/op",
            "extra": "164880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - ns/op",
            "value": 7281,
            "unit": "ns/op",
            "extra": "164880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - B/op",
            "value": 5842,
            "unit": "B/op",
            "extra": "164880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/filter - allocs/op",
            "value": 105,
            "unit": "allocs/op",
            "extra": "164880 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map",
            "value": 5956,
            "unit": "ns/op\t    4665 B/op\t      84 allocs/op",
            "extra": "200499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - ns/op",
            "value": 5956,
            "unit": "ns/op",
            "extra": "200499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - B/op",
            "value": 4665,
            "unit": "B/op",
            "extra": "200499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertComprehensions/map - allocs/op",
            "value": 84,
            "unit": "allocs/op",
            "extra": "200499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access",
            "value": 3491,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "335214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - ns/op",
            "value": 3491,
            "unit": "ns/op",
            "extra": "335214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "335214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/simple_access - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "335214 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access",
            "value": 4135,
            "unit": "ns/op\t    2264 B/op\t      40 allocs/op",
            "extra": "291369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - ns/op",
            "value": 4135,
            "unit": "ns/op",
            "extra": "291369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - B/op",
            "value": 2264,
            "unit": "B/op",
            "extra": "291369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_access - allocs/op",
            "value": 40,
            "unit": "allocs/op",
            "extra": "291369 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has",
            "value": 2778,
            "unit": "ns/op\t    1664 B/op\t      27 allocs/op",
            "extra": "430296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - ns/op",
            "value": 2778,
            "unit": "ns/op",
            "extra": "430296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - B/op",
            "value": 1664,
            "unit": "B/op",
            "extra": "430296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_has - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "430296 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has",
            "value": 3227,
            "unit": "ns/op\t    1816 B/op\t      30 allocs/op",
            "extra": "368547 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - ns/op",
            "value": 3227,
            "unit": "ns/op",
            "extra": "368547 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - B/op",
            "value": 1816,
            "unit": "B/op",
            "extra": "368547 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/nested_json_has - allocs/op",
            "value": 30,
            "unit": "allocs/op",
            "extra": "368547 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison",
            "value": 3484,
            "unit": "ns/op\t    2024 B/op\t      35 allocs/op",
            "extra": "342926 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - ns/op",
            "value": 3484,
            "unit": "ns/op",
            "extra": "342926 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - B/op",
            "value": 2024,
            "unit": "B/op",
            "extra": "342926 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/json_comparison - allocs/op",
            "value": 35,
            "unit": "allocs/op",
            "extra": "342926 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json",
            "value": 7448,
            "unit": "ns/op\t    3785 B/op\t      64 allocs/op",
            "extra": "160822 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - ns/op",
            "value": 7448,
            "unit": "ns/op",
            "extra": "160822 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - B/op",
            "value": 3785,
            "unit": "B/op",
            "extra": "160822 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertJSONPath/complex_json - allocs/op",
            "value": 64,
            "unit": "allocs/op",
            "extra": "160822 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern",
            "value": 8284,
            "unit": "ns/op\t    5981 B/op\t      75 allocs/op",
            "extra": "143422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - ns/op",
            "value": 8284,
            "unit": "ns/op",
            "extra": "143422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - B/op",
            "value": 5981,
            "unit": "B/op",
            "extra": "143422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/simple_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "143422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive",
            "value": 7920,
            "unit": "ns/op\t    5975 B/op\t      75 allocs/op",
            "extra": "151725 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - ns/op",
            "value": 7920,
            "unit": "ns/op",
            "extra": "151725 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - B/op",
            "value": 5975,
            "unit": "B/op",
            "extra": "151725 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/case_insensitive - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "151725 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern",
            "value": 8415,
            "unit": "ns/op\t    5983 B/op\t      75 allocs/op",
            "extra": "142968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - ns/op",
            "value": 8415,
            "unit": "ns/op",
            "extra": "142968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - B/op",
            "value": 5983,
            "unit": "B/op",
            "extra": "142968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/complex_pattern - allocs/op",
            "value": 75,
            "unit": "allocs/op",
            "extra": "142968 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class",
            "value": 8076,
            "unit": "ns/op\t    5992 B/op\t      76 allocs/op",
            "extra": "147073 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - ns/op",
            "value": 8076,
            "unit": "ns/op",
            "extra": "147073 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - B/op",
            "value": 5992,
            "unit": "B/op",
            "extra": "147073 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_digit_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "147073 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class",
            "value": 7951,
            "unit": "ns/op\t    5993 B/op\t      76 allocs/op",
            "extra": "148965 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - ns/op",
            "value": 7951,
            "unit": "ns/op",
            "extra": "148965 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - B/op",
            "value": 5993,
            "unit": "B/op",
            "extra": "148965 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_class - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "148965 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary",
            "value": 8153,
            "unit": "ns/op\t    5985 B/op\t      76 allocs/op",
            "extra": "145159 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - ns/op",
            "value": 8153,
            "unit": "ns/op",
            "extra": "145159 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - B/op",
            "value": 5985,
            "unit": "B/op",
            "extra": "145159 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertRegex/with_word_boundary - allocs/op",
            "value": 76,
            "unit": "allocs/op",
            "extra": "145159 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5",
            "value": 9875,
            "unit": "ns/op\t    6578 B/op\t     106 allocs/op",
            "extra": "121849 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - ns/op",
            "value": 9875,
            "unit": "ns/op",
            "extra": "121849 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - B/op",
            "value": 6578,
            "unit": "B/op",
            "extra": "121849 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_5 - allocs/op",
            "value": 106,
            "unit": "allocs/op",
            "extra": "121849 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10",
            "value": 19108,
            "unit": "ns/op\t   13404 B/op\t     197 allocs/op",
            "extra": "62758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - ns/op",
            "value": 19108,
            "unit": "ns/op",
            "extra": "62758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - B/op",
            "value": 13404,
            "unit": "B/op",
            "extra": "62758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_and_10 - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "62758 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5",
            "value": 2116,
            "unit": "ns/op\t    1592 B/op\t      25 allocs/op",
            "extra": "564205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - ns/op",
            "value": 2116,
            "unit": "ns/op",
            "extra": "564205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - B/op",
            "value": 1592,
            "unit": "B/op",
            "extra": "564205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_parentheses_5 - allocs/op",
            "value": 25,
            "unit": "allocs/op",
            "extra": "564205 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary",
            "value": 8804,
            "unit": "ns/op\t    6122 B/op\t     100 allocs/op",
            "extra": "136576 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - ns/op",
            "value": 8804,
            "unit": "ns/op",
            "extra": "136576 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - B/op",
            "value": 6122,
            "unit": "B/op",
            "extra": "136576 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_ternary - allocs/op",
            "value": 100,
            "unit": "allocs/op",
            "extra": "136576 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic",
            "value": 5471,
            "unit": "ns/op\t    3705 B/op\t      67 allocs/op",
            "extra": "221038 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - ns/op",
            "value": 5471,
            "unit": "ns/op",
            "extra": "221038 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - B/op",
            "value": 3705,
            "unit": "B/op",
            "extra": "221038 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertDeeplyNested/nested_arithmetic - allocs/op",
            "value": 67,
            "unit": "allocs/op",
            "extra": "221038 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20",
            "value": 42827,
            "unit": "ns/op\t   29106 B/op\t     424 allocs/op",
            "extra": "27702 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - ns/op",
            "value": 42827,
            "unit": "ns/op",
            "extra": "27702 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - B/op",
            "value": 29106,
            "unit": "B/op",
            "extra": "27702 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/mixed_conditions_20 - allocs/op",
            "value": 424,
            "unit": "allocs/op",
            "extra": "27702 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain",
            "value": 39452,
            "unit": "ns/op\t   27450 B/op\t     398 allocs/op",
            "extra": "30316 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - ns/op",
            "value": 39452,
            "unit": "ns/op",
            "extra": "30316 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - B/op",
            "value": 27450,
            "unit": "B/op",
            "extra": "30316 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertLargeExpression/large_and_chain - allocs/op",
            "value": 398,
            "unit": "allocs/op",
            "extra": "30316 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison",
            "value": 3796,
            "unit": "ns/op\t    2528 B/op\t      46 allocs/op",
            "extra": "311422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - ns/op",
            "value": 3796,
            "unit": "ns/op",
            "extra": "311422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - B/op",
            "value": 2528,
            "unit": "B/op",
            "extra": "311422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/timestamp_comparison - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "311422 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function",
            "value": 1822,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "662476 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - ns/op",
            "value": 1822,
            "unit": "ns/op",
            "extra": "662476 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "662476 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/date_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "662476 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function",
            "value": 1848,
            "unit": "ns/op\t    1512 B/op\t      27 allocs/op",
            "extra": "639950 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - ns/op",
            "value": 1848,
            "unit": "ns/op",
            "extra": "639950 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - B/op",
            "value": 1512,
            "unit": "B/op",
            "extra": "639950 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertTimestamps/datetime_function - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "639950 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith",
            "value": 2212,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "524499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - ns/op",
            "value": 2212,
            "unit": "ns/op",
            "extra": "524499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "524499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/startsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "524499 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith",
            "value": 2275,
            "unit": "ns/op\t    1656 B/op\t      28 allocs/op",
            "extra": "519841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - ns/op",
            "value": 2275,
            "unit": "ns/op",
            "extra": "519841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - B/op",
            "value": 1656,
            "unit": "B/op",
            "extra": "519841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/endsWith - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "519841 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains",
            "value": 2172,
            "unit": "ns/op\t    1648 B/op\t      27 allocs/op",
            "extra": "533253 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - ns/op",
            "value": 2172,
            "unit": "ns/op",
            "extra": "533253 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - B/op",
            "value": 1648,
            "unit": "B/op",
            "extra": "533253 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/contains - allocs/op",
            "value": 27,
            "unit": "allocs/op",
            "extra": "533253 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation",
            "value": 3186,
            "unit": "ns/op\t    2176 B/op\t      37 allocs/op",
            "extra": "372244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - ns/op",
            "value": 3186,
            "unit": "ns/op",
            "extra": "372244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - B/op",
            "value": 2176,
            "unit": "B/op",
            "extra": "372244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/concatenation - allocs/op",
            "value": 37,
            "unit": "allocs/op",
            "extra": "372244 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops",
            "value": 6155,
            "unit": "ns/op\t    4121 B/op\t      72 allocs/op",
            "extra": "193314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - ns/op",
            "value": 6155,
            "unit": "ns/op",
            "extra": "193314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - B/op",
            "value": 4121,
            "unit": "B/op",
            "extra": "193314 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertStringOperations/multiple_string_ops - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193314 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison",
            "value": 6181,
            "unit": "ns/op\t    4083 B/op\t      72 allocs/op",
            "extra": "193329 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - ns/op",
            "value": 6181,
            "unit": "ns/op",
            "extra": "193329 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - B/op",
            "value": 4083,
            "unit": "B/op",
            "extra": "193329 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/simple_comparison - allocs/op",
            "value": 72,
            "unit": "allocs/op",
            "extra": "193329 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path",
            "value": 6170,
            "unit": "ns/op\t    4025 B/op\t      70 allocs/op",
            "extra": "193146 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - ns/op",
            "value": 6170,
            "unit": "ns/op",
            "extra": "193146 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - B/op",
            "value": 4025,
            "unit": "B/op",
            "extra": "193146 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/json_path - allocs/op",
            "value": 70,
            "unit": "allocs/op",
            "extra": "193146 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern",
            "value": 12329,
            "unit": "ns/op\t    8319 B/op\t     117 allocs/op",
            "extra": "99093 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - ns/op",
            "value": 12329,
            "unit": "ns/op",
            "extra": "99093 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - B/op",
            "value": 8319,
            "unit": "B/op",
            "extra": "99093 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/regex_pattern - allocs/op",
            "value": 117,
            "unit": "allocs/op",
            "extra": "99093 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation",
            "value": 6890,
            "unit": "ns/op\t    4403 B/op\t      81 allocs/op",
            "extra": "171360 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - ns/op",
            "value": 6890,
            "unit": "ns/op",
            "extra": "171360 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - B/op",
            "value": 4403,
            "unit": "B/op",
            "extra": "171360 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/array_operation - allocs/op",
            "value": 81,
            "unit": "allocs/op",
            "extra": "171360 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query",
            "value": 22634,
            "unit": "ns/op\t   13938 B/op\t     217 allocs/op",
            "extra": "52908 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - ns/op",
            "value": 22634,
            "unit": "ns/op",
            "extra": "52908 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - B/op",
            "value": 13938,
            "unit": "B/op",
            "extra": "52908 times\n4 procs"
          },
          {
            "name": "BenchmarkAnalyzeQuery/complex_query - allocs/op",
            "value": 217,
            "unit": "allocs/op",
            "extra": "52908 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options",
            "value": 2806,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "419036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - ns/op",
            "value": 2806,
            "unit": "ns/op",
            "extra": "419036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "419036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/no_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "419036 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas",
            "value": 2827,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "418237 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - ns/op",
            "value": 2827,
            "unit": "ns/op",
            "extra": "418237 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "418237 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_schemas - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "418237 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth",
            "value": 2815,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "423544 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - ns/op",
            "value": 2815,
            "unit": "ns/op",
            "extra": "423544 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "423544 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_depth - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "423544 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output",
            "value": 2811,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "421714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - ns/op",
            "value": 2811,
            "unit": "ns/op",
            "extra": "421714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "421714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/with_max_output - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "421714 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options",
            "value": 2829,
            "unit": "ns/op\t    1840 B/op\t      31 allocs/op",
            "extra": "422656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - ns/op",
            "value": 2829,
            "unit": "ns/op",
            "extra": "422656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - B/op",
            "value": 1840,
            "unit": "B/op",
            "extra": "422656 times\n4 procs"
          },
          {
            "name": "BenchmarkConvertWithOptions/all_options - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "422656 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized",
            "value": 627509,
            "unit": "ns/op\t   14272 B/op\t    1713 allocs/op",
            "extra": "1914 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - ns/op",
            "value": 627509,
            "unit": "ns/op",
            "extra": "1914 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - B/op",
            "value": 14272,
            "unit": "B/op",
            "extra": "1914 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Parameterized - allocs/op",
            "value": 1713,
            "unit": "allocs/op",
            "extra": "1914 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline",
            "value": 522271,
            "unit": "ns/op\t   14040 B/op\t    1708 allocs/op",
            "extra": "2294 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - ns/op",
            "value": 522271,
            "unit": "ns/op",
            "extra": "2294 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - B/op",
            "value": 14040,
            "unit": "B/op",
            "extra": "2294 times\n4 procs"
          },
          {
            "name": "BenchmarkParameterizedVsInline/Inline - allocs/op",
            "value": 1708,
            "unit": "allocs/op",
            "extra": "2294 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small",
            "value": 602.7,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1984940 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - ns/op",
            "value": 602.7,
            "unit": "ns/op",
            "extra": "1984940 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1984940 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Small - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1984940 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium",
            "value": 606.2,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1984268 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - ns/op",
            "value": 606.2,
            "unit": "ns/op",
            "extra": "1984268 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1984268 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Medium - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1984268 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large",
            "value": 608.4,
            "unit": "ns/op\t     112 B/op\t       5 allocs/op",
            "extra": "1971001 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - ns/op",
            "value": 608.4,
            "unit": "ns/op",
            "extra": "1971001 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - B/op",
            "value": 112,
            "unit": "B/op",
            "extra": "1971001 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldLookup_Large - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1971001 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small",
            "value": 144.4,
            "unit": "ns/op\t     176 B/op\t       2 allocs/op",
            "extra": "8325703 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - ns/op",
            "value": 144.4,
            "unit": "ns/op",
            "extra": "8325703 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - B/op",
            "value": 176,
            "unit": "B/op",
            "extra": "8325703 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Small - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8325703 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large",
            "value": 9092,
            "unit": "ns/op\t   16400 B/op\t       2 allocs/op",
            "extra": "150578 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - ns/op",
            "value": 9092,
            "unit": "ns/op",
            "extra": "150578 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - B/op",
            "value": 16400,
            "unit": "B/op",
            "extra": "150578 times\n4 procs"
          },
          {
            "name": "BenchmarkFieldNames_Large - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "150578 times\n4 procs"
          }
        ]
      }
    ]
  }
}