# markup

_module_
A generalized data format to represent human readable contents.
Though it is not meant to represent HTML nodes.

- _data_ [Bold](#Bold)
- _data_ [Code](#Code)
- _data_ [CodeBlock](#CodeBlock)
- _data_ [Format](#Format)
- _enum_ [FormatWitness](#FormatWitness)
- _data_ [Heading](#Heading)
- _data_ [Image](#Image)
- _data_ [Italic](#Italic)
- _data_ [Link](#Link)
- _enum_ [Markup](#Markup)
- _data_ [OrderedList](#OrderedList)
- _data_ [Paragraph](#Paragraph)
- _data_ [UnorderedList](#UnorderedList)
- _func_ [b](#b)
- _func_ [code](#code)
- _func_ [convert](#convert) markup, targetFormatWitness
- _func_ [from](#from) witness
- _func_ [group](#group) list
- _func_ [h1](#h1)
- _func_ [h2](#h2)
- _func_ [h3](#h3)
- _func_ [h4](#h4)
- _func_ [h5](#h5)
- _func_ [h6](#h6)
- _func_ [i](#i)
- _func_ [img](#img)
- _func_ [inline](#inline) list
- _func_ [li](#li) item
- _func_ [link](#link)
- _func_ [ol](#ol) child
- _func_ [p](#p) child
- _func_ [pre](#pre)
- _func_ [ref](#ref) child, name
- _func_ [text](#text) str
- _func_ [ul](#ul) child

## Bold

_data_ A bold portion of inline text.

> Bold **"text"**

### Properties

- `child` - At least one child of type Markup.

## Code

_data_ Some preformatted inline code.

> `Code text`

### Properties

- `text` - Some text. Markup not supported.

## CodeBlock

_data_ A block of preformatted multiline code.

> ```lithia
> CodeBlock "lithia", "text"
> ```

### Properties

- `language` - Maybe the language as string.
Implementations must support Some, None and String.
- `text` - The text itself as String.

## Format

_data_ A specific markup format to convert to.

### Properties

- `convert markup` - Converts given markup to this format.

## FormatWitness

_enum_
A more generalized form of a markup.Format.
Allows to use plain functions and whole modules as format.

### Cases

- [Format](#Format)
- [Module](#Module)
- [Function](#Function)

## Heading

_data_ A heading of a specific nesting level.

> ##### Heading 5, "title"

### Properties

- `level` - The level of nesting.
An Int, typically in the range of one to six.
- `child` - At least one child of type Markup.

## Image

_data_ An embedded image.

> Image ![url](https://github.com/vknabel/lithia/blob/main/assets/lithia.png), alt

### Properties

- `url` - The url  of the image.
Markup not supported.
- `alt` - Alternatively displayed text in case the image can't be loaded or for accessibility.
Markup not supported.

## Italic

_data_ An italic portion of inline text.

> Italic _"text"_

### Properties

- `child` - At least one child of type Markup.

## Link

_data_ A link to a resource.

> [Link child, url](https://github.com/vknabel/lithia)

### Properties

- `child` - At least one child of type Markup.
Used to display the link.
- `url` - The url to point to.
Markup not supported.

## Markup

_enum_
A recursive enum of markup nodes.

### Cases

- [String](#String)
- [List](#List)
- [Heading](#Heading)
- [Paragraph](#Paragraph)
- [Italic](#Italic)
- [Bold](#Bold)
- [Link](#Link)
- [Image](#Image)
- [Code](#Code)
- [CodeBlock](#CodeBlock)
- [UnorderedList](#UnorderedList)
- [OrderedList](#OrderedList)

## OrderedList

_data_ An ordered list. Requires a list of children.

> OrderedList [
>   1. item
>   2. another
>
> ]

### Properties

- `children` - A list of Markup children.

## Paragraph

_data_ A paragraph with some blank lines around.

> Paragraph [
>     "First",
>     "Second"
> ]

### Properties

- `child` - At least one child of type Markup.

## UnorderedList

_data_ An unordered list. Requires a list of children.

> UnorderedList [
>   - item
>   - another
>
> ]

### Properties

- `children` - A list of Markup children.

## b

_func_ `b`

**bold** text

## code

_func_ `code`

A Code block.

## convert

_func_ `convert markup, targetFormatWitness`

Converts some markup to a target format.

## from

_func_ `from witness`

Creates a Format from a given FormatWitness.

## group

_func_ `group list`

A list of markup for better readbility.

## h1

_func_ `h1`

## h2

_func_ `h2`

## h3

_func_ `h3`

## h4

_func_ `h4`

## h5

_func_ `h5`

## h6

_func_ `h6`

## i

_func_ `i`

_italic_ text

## img

_func_ `img`

An image with url and alt.

## inline

_func_ `inline list`

Some inline text. Ignores empty elements. Adds spaces.

## li

_func_ `li item`

A list item for better readability.

## link

_func_ `link`

A link with title and url.

## ol

_func_ `ol child`

An ordered list. Ignores empty lines.

## p

_func_ `p child`

A paragraph. Ignores empty lines.

## pre

_func_ `pre`

Preformatted inline code.

## ref

_func_ `ref child, name`

An internal ref to a heading.

## text

_func_ `text str`

Constant text for better readbility.

## ul

_func_ `ul child`

An unordered list. Ignores empty lines.

