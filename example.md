It's time to continue exploring the fundamentals of I/O in Go.

If you’re not familiar with io.Reader and io.Writer in Go, read [my first article](/@andreiboar/fundamentals-of-i-o-in-go-c893d3714deb), especially the part where I said that a Reader is something we can read from, and a Writer is something we can write to.

Let's see what else is out there and show you some awesome functions.

# io.LimitReader

Sometimes, you want to limit the number of bytes you read from a piece of data.

To limit the size of data you read, you can wrap that data source (Reader) with an [io.LimitReader](https://pkg.go.dev/io#LimitReader) function that returns another Reader from which you can read only `n` bytes:

```
package main

import (
    "io"
    "log"
    "os"
    "strings"
)

func main() {
    r := strings.NewReader("Hello!")

    lr := io.LimitReader(r, 4)

    // Ouput: Hell
    if _, err := io.Copy(os.Stdout, lr); err != nil {
        log.Fatal(err)
    }
}
```

[run in Playground](https://go.dev/play/p/1GYziqBeupc)

The above code outputs `Hell`. The limited reader `lr` can be passed into any function expecting a Reader and will work as expected: only 4 bytes will be read from it.

# io.MultiReader

Sometimes, you can have multiple data sources, and you want to treat those sources as one. If you have multiple Readers, you can merge them into one Reader with [io.MultiReader](https://pkg.go.dev/io#MultiReader):

```
package main

import (
    "io"
    "log"
    "os"
    "strings"
)

func main() {
    r1 := strings.NewReader("first reader\n")
    r2 := strings.NewReader("second reader\n")
    r3 := strings.NewReader("third reader\n")

    // merge all 3 readers into one
    r := io.MultiReader(r1, r2, r3)

    if _, err := io.Copy(os.Stdout, r); err != nil {
       log.Fatal(err)
    }

}
```

[run in Playground](https://go.dev/play/p/qrXFNukv11n)

In the above example, we merge all three readers into `r`, and when we read to `os.Stdout` from `r`, it's as if we copy sequentially from each reader.

I recently used io.MultiReader, when I needed to write the bytes coming from different requests into one file in a [parallel downloader implementation](https://github.com/zuzuleinen/downloader/blob/main/downloader/downloader.go#L118):

```
 results := make([]io.Reader, n)
 // ...
 if err := writeToFile(destinationFileName, io.MultiReader(results...)); err != nil {
  return fmt.Errorf("could not write to file: %w", err)
 }
```

# io.MultiWriter

Similar to io.MultiReader, we have an [io.MultiWriter](https://pkg.go.dev/io#MultiWriter) function, which creates a writer that reproduces its writes to all the provided writers:

```
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

func main() {
    var (
       buf1 bytes.Buffer
       buf2 bytes.Buffer
    )

    // writing into mw will write to both buf1 and buf2
    mw := io.MultiWriter(&buf1, &buf2)

    // r is the source of data(Reader)
    r := strings.NewReader("some io.Reader stream to be read")

    // write to mw from r
    io.Copy(mw, r)

    fmt.Println("data inside buffer1 :", buf1.String())
    fmt.Println("data inside buffer2 :", buf2.String())

}
```

[run in Playround](https://go.dev/play/p/pgGn2NO95Qc)

I like to use `io.MultiWriter` when I'm trying to debug what was written in a certain Writer.

If a function writes to a Writer and, for some reason, it is too hard to get those contents, I connect a `bytes.Buffer` to it, and then I check the contents of the buffer, which will be the same as the contents of my inaccessible Writer:

```
package main

import (
    "bytes"
    "fmt"
    "io"
)

func main() {
    buf := new(bytes.Buffer)

    debug := io.MultiWriter(buf, weirdWriter) // attach buf to weirdWriter

    complicatedFunctionWithAWriter(weirdWriter)

    // The contents of the buffer will be the same as in weirdWriter
    fmt.Println(buf.String())
    fmt.Println(buf.Bytes())

}
```

# io.TeeReader

Imagine reading and writing data in a place, and you want to write the same data somewhere else.

In the above image, we read from `R` to ` W` and write to an extra `Logs` Writer at the same time.

Let’s see how that looks in code:

```
package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    logs := new(bytes.Buffer)

    data := strings.NewReader("Hello World!\n")

    teeReader := io.TeeReader(data, logs)

    // logs will also receives contents from teeReader
    io.Copy(os.Stdout, teeReader)

    fmt.Println("Content of logs:", logs.String())
}
```

[run in Playground](https://go.dev/play/p/tJmeNhcrSXt)

The behavior of the TeeReader is exactly like the [tee](https://www.geeksforgeeks.org/tee-command-linux-example/) command from Linux, named after the T-splitter used in plumbing:

# io.Pipe

Speaking of Linux and plumbing, you might be familiar with the [pipe operator](https://www.geeksforgeeks.org/piping-in-unix-or-linux/), which combines two or more commands so that the output of one becomes the input of the other.

```
echo hello | tr l y
```

This outputs:

```
heyyo
```

I like the pipe's universality. It doesn't matter what programs you connect as long as one writes and the other reads. I could have just as easily used `cat` to fetch the contents of a file instead of `echo`:

```
cat file | tr l y
```

In Go, we achieve this behavior with [io.Pipe](https://pkg.go.dev/io#Pipe) that can be used to connect code expecting an io.Reader with code expecting an io.Writer.

Let's try to replicate the same `echo hello | tr l y` example to see how that works:

```
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    pipeReader, pipeWriter := io.Pipe()

    echo(pipeWriter, "hello")
    tr(pipeReader, "e", "i")
}

func echo(w io.Writer, s string) {
    fmt.Fprint(w, s)
}

func tr(r io.Reader, old string, new string) {
    data, _ := io.ReadAll(r)
    res := strings.Replace(string(data), old, new, -1)
    fmt.Println(res)
}

```

[run in Playground](https://go.dev/play/p/df0TKZUtuoX)

Running this program, we get an error:

```
fatal error: all goroutines are asleep - deadlock!
```

To understand why, we need to read the documentation of the `io.Pipe` function:

```
// Pipe creates a synchronous in-memory pipe.
// That is, each Write to the [PipeWriter] blocks until it has satisfied
// one or more Reads from the [PipeReader] that fully consume
// the written data.
```

Our code is not synchronous but sequential: we call `echo` and then `tr`. We get a deadlock when we try to write because no reading is happening.

Let's fix that:

```
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    pipeReader, pipeWriter := io.Pipe()

    // Run echo concurrently with tr in a separate goroutine
    go echo(pipeWriter, "hello")
    tr(pipeReader, "e", "i")
}

func echo(w io.Writer, s string) {
    fmt.Fprint(w, s)
}

func tr(r io.Reader, old string, new string) {
    data, _ := io.ReadAll(r)
    res := strings.Replace(string(data), old, new, -1)
    fmt.Println(res)
}
```

[run in Playground](https://go.dev/play/p/vA-6oIxJJv9)

Running the modified program gives us the same error:

```
fatal error: all goroutines are asleep - deadlock!
```

Why is that? Let's check the `PipeReader` documentation on the Read method:

```
// Read implements the standard Read interface:
// it reads data from the pipe, blocking until a writer
// arrives or the write end is closed.
// If the write end is closed with an error, that error is
// returned as err; otherwise err is EOF.
func (r *PipeReader) Read(data []byte) (n int, err error) {
    return r.pipe.read(data)
}
```

The second part of this sentence is key: blocking until a writer arrives or the write-end **is closed**.

The read is still blocked because the writer is not closed after writing our data. Let's fix that by calling `Close()` on the `pipeWriter`:

```
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    pipeReader, pipeWriter := io.Pipe()

    go func() {
       echo(pipeWriter, "hello")
       // we close the writer so we unblock the reader
       pipeWriter.Close()
    }()
    tr(pipeReader, "e", "i")
}

func echo(w io.Writer, s string) {
    fmt.Fprint(w, s)
}

func tr(r io.Reader, old string, new string) {
    data, _ := io.ReadAll(r)
    res := strings.Replace(string(data), old, new, -1)
    fmt.Println(res)
}
```

[run in Playground](https://go.dev/play/p/ngjj4_bMt0h)

## Conclusion

That was all for today! I hope you found this post helpful and that you now have a couple more functions under your belt.

**Covered functions:**

[io.LimitReader](https://pkg.go.dev/io#LimitReader)

[io.MultiReader](https://pkg.go.dev/io#MultiReader)

[io.MultiWriter](https://pkg.go.dev/io#MultiWriter)

[io.TeeReader](https://pkg.go.dev/io#TeeReader)

[io.Pipe](https://pkg.go.dev/io#Pipe)

If you enjoyed this article, I would appreciate a clap or share!

Do you have any questions? Connect with me on [LinkedIn](https://www.linkedin.com/in/andrei-boar/), and I'll be happy to help!