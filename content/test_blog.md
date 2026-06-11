---
Title: the-test
Tags: test blog
Date: 2026-06-08
Desc: A test blog
---

This is a test blog entry

# This is a header
foo blah

## This is a smaller header?
baz

- this
- is 
- a list
    - that 
        - I 
            - am testing
- foo

1. Numbered list
1. Numbered list
2. Numbered list

[link](https://dylanmeeks.engineer)
[link2][1]

> Quote
> line 2
>> second level
> back to first

---

![Image](http://url/a.png)

---

This is `inline code`.

```zig
const std = @import("std");

fn main(init: std.Init) !void {
    const io = init.io;
    std.debug.print("This is a test of highlighting and code embed", .{});
}
```

[1]: https://dylanmeeks.engineer/blog
