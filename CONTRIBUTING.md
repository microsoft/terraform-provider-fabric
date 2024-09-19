# 👥 Contributing Guide

This project welcomes feedback and suggestions only via GitHub Issues. Pull Request (PR) contributions will **NOT** be accepted at this time. In the future, PRs may be accepted, in which case you will be require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution. For details, visit <https://cla.opensource.microsoft.com>.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## 🔄️ Contribution process (Microsoft FTE only)

To contribute, please follow these steps:

1. Fork/clone the project repository on GitHub.
1. Create a new branch for your feature or bug fix.
1. Make sure the `README.md` and any other relevant documentation are kept up-to-date.
1. Make your changes and commit them with descriptive commit messages; check [Conventional Commits](https://www.conventionalcommits.org) as a suggestion.
1. Push to your forked repository or branch.
1. Create a new Pull Request from your fork/branch to this project.
1. Please ensure that your pull request includes a detailed description of your changes and that your code adheres to the code style guidelines outlined below.

## ✍️ Types of Contributions

We welcome feedback and suggestions only via GitHub Issues. Pull Request (PR) contributions will **NOT** be accepted at this time - it may change in the future.

### Resources

Creating a new [resource](https://developer.hashicorp.com/terraform/plugin/framework/resources) can allow terraform to manage new infrastructure/services not currently provided by the provider.

### Data Sources

Creating a new [data source](https://developer.hashicorp.com/terraform/plugin/framework/data-sources) can allow terraform to reference data about infrastructure and services.

### Examples

Examples of real-world use cases are encouraged. Please contribute those types of examples to the [Fabric Terraform QuickStarts](https://aka.ms/FabricTF/quickstart) repo.

## ☑️ Pull Request Checklist

PRs for new resources or data sources are expected to meeting the following criteria:

- Add a production quality implementation of the resource or data-source in [./internal/provider/services](./internal/provider/services)
- Add unit tests and acceptance tests for your contribution in [./internal/provider](./internal/provider)
  - Tests should pass and provide >90% coverage of your contribution
- Add examples for your contribution in [./examples](./examples) (see [Terraform Documentation on examples](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-configuration-examples))
- Add [schema descriptions](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-schema-descriptions) for your resource or data-source in [./internal/provider/services](./internal/provider/services)
- and/or [./templates](./templates)
- Update auto-generated documentation in [./docs](./docs). (DO NOT manually edit [./docs](./docs) or your updates will be overwritten)
- Ensure the PR description clearly describes the feature you're adding and any known limitations

## 🔰 Code of Conduct

All contributors are expected to adhere to the project name code of conduct. Therefore, please review it before contributing [`Code of Conduct`](./CODE_OF_CONDUCT.md).

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the project license.

---

Thank you for contributing!
