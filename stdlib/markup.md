# markup

_module_

- _data_ [Bold](#Bold)
- _data_ [Code](#Code)
- _data_ [CodeBlock](#CodeBlock)
- _data_ [Heading](#Heading)
- _data_ [Image](#Image)
- _data_ [Italic](#Italic)
- _data_ [Link](#Link)
- _enum_ [MarkupNode](#MarkupNode)
- _data_ [OrderedList](#OrderedList)
- _data_ [Paragraph](#Paragraph)
- _data_ [Serializer](#Serializer)
- _enum_ [SerializerWitness](#SerializerWitness)
- _data_ [UnorderedList](#UnorderedList)
- _func_ [b](#b)
- _func_ [code](#code)
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
- _func_ [serialize](#serialize) markup, witness
- _func_ [text](#text) str
- _func_ [ul](#ul) child

## Bold

_data_

### Properties

- `child`

## Code

_data_

### Properties

- `text`

## CodeBlock

_data_

### Properties

- `language`
- `text`

## Heading

_data_

### Properties

- `level`
- `child`

## Image

_data_

### Properties

- `url`
- `alt`

## Italic

_data_

### Properties

- `child`

## Link

_data_

### Properties

- `child`
- `url`

## MarkupNode

_enum_

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

_data_

### Properties

- `children`

## Paragraph

_data_

### Properties

- `child`

## Serializer

_data_

### Properties

- `serialize markup`

## SerializerWitness

_enum_

### Cases

- [Serializer](#Serializer)
- [Module](#Module)
- [Function](#Function)

## UnorderedList

_data_

### Properties

- `children`

## b

_func_ `b`

## code

_func_ `code`

## from

_func_ `from witness`

## group

_func_ `group list`

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

## img

_func_ `img`

## inline

_func_ `inline list`

## li

_func_ `li item`

## link

_func_ `link`

## ol

_func_ `ol child`

## p

_func_ `p child`

## pre

_func_ `pre`

## ref

_func_ `ref child, name`

## serialize

_func_ `serialize markup, witness`

## text

_func_ `text str`

## ul

_func_ `ul child`

