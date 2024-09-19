def on_page_markdown(markdown, **kwargs):
    # https://developer.hashicorp.com/terraform/registry/providers/docs#callouts
    # https://squidfunk.github.io/mkdocs-material/reference/admonitions/
    markdown = markdown.replace("-> ", "!!! note \"Note\"\n    ")
    markdown = markdown.replace("~> ", "!!! warning \"Note\"\n    ")
    markdown = markdown.replace("!> ", "!!! danger \"Warning\"\n    ")

    markdown = markdown.replace(" (Function)", "")
    markdown = markdown.replace(" (Resource)", "")
    markdown = markdown.replace(" (Data Source)", "")

    return markdown
