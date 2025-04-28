# Package yconf

一个使用YMAL风格语言的配置文件管理。

为了与jconf包保持结果和逻辑一致，用了"github.com/goccy/go-yaml"将YAML转成了JSON后，仍然利用golang官方的"encoding/json"进行后续的处理。