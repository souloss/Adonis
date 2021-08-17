# Adonis

## Abstract
Adonis 是一个针对 Kubernetes 资源清单设计的简化镜像操作的工具

主要用于从文件或文件夹中的资源清单（YAML 格式）中搜寻使用的镜像文件，对其进行统一操作，如 pull、tag、save、push 等

Adonis 需要借助 Dockerd 来完成镜像的相关操作

## Usage
Adonis 有 5 个子命令：parse、pull、save、tag 和 push

有全局标志 -f / --files，用于指定需要操作的文件或者文件夹，如果指定文件夹，则会对该文件夹下所有合法的 YAML 文件进行操作

### parse
该命令会从文件中获取镜像
```shell
adonis parse -f <file_path/dir_path>
```

### pull
该命令会找到文件中使用的所有镜像，并拉取
```shell
adonis pull -f <file_path/dir_path>
```

### save
在拉取镜像的基础之上，将镜像保存到指定路径（默认为 ./）

文件名称格式为 镜像仓库名称的最后一段 + _ + 标签名称 + .tar

如 `eipwork/etcd-host:3.4.16-1` 会保存为 `etcd-host_3.4.16-1.tar`

```shell
adonis pull -f <file_path/dir_path>
```

还可以使用 -p / --path 指定镜像文件保存的目录，如果目录不存在会自行创建

```shell
adonis pull -f <file_path/dir_path> -p <save_path>
```

### tag
在拉取镜像的基础之上，为拉取到的镜像打标签

需要使用 -r / --repo 指定新的镜像仓库地址
```shell
adonis tag -f <file_path/dir_path> -r <new_repository_path>
```

还可以使用 -d / --delete 参数，删除原有的标签
```shell
adonis tag -f <file_path/dir_path> -r <new_repository_path> -d
```

### push
在为镜像打标签的基础之上，将其推送到新标签所对应的仓库

```shell
adonis push -f <file_path/dir_path> -r <new_repository_path>
```

和 save 命令一样，可以通过 -s / --save 保存镜像，使用 -p / --path 指定保存路径（默认为 ./）

类似于 tag 命令，可以使用 -d / --delete 删除原来的镜像标签
