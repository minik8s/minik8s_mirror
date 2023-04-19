// 测试ImageManager的相关函数
package image

import "testing"

var TestImageURLs = []string{
	"docker.io/library/busybox:latest",
	"docker.io/library/hello-world:latest",
	"docker.io/library/nginx:latest",
	"docker.io/library/redis:latest",
}

// 删除所有存在的镜像
// 这个函数在所有的测试函数执行之前执行
func TestMain(m *testing.M) {
	// 创建一个ImageManager
	im := &ImageManager{}
	// 获取所有已经存在的镜像
	images, err := im.ListAllImages()
	if err != nil {
		panic(err)
	}

	// 遍历所有的镜像，删除它们
	for _, image := range images {
		err := im.RemoveImage(image.Name)
		if err != nil {
			panic(err)
		}
	}
	m.Run()
}

// 测试PullImage函数
func TestPullImage(t *testing.T) {
	// 创建一个ImageManager
	im := &ImageManager{}

	// 遍历imageURLs
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImage(imageURL)

		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}
}

// 测试GetImage函数
func TestGetImage(t *testing.T) {
	// 创建一个ImageManager，对应的是接口
	im := &ImageManager{}
	// 遍历imageURLs
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImage(imageURL)

		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
		if image.Name() != imageURL {
			t.Errorf("PullImage() failed: %s\n", image.Name())
		}
	}

}

// 测试ListAllImages函数
func TestListAllImages(t *testing.T) {
	// 创建一个ImageManager
	im := &ImageManager{}
	// 调用ListAllImages函数
	images, err := im.ListAllImages()
	if err != nil {
		t.Errorf("ListAllImages() failed: %s\n", err)
	}
	if images == nil {
		t.Errorf("ListAllImages() failed: images is null")
	}
	// 遍历images
	for index, image := range images {
		nameStr := image.Name
		// 检查是否和预期的一致
		if nameStr != TestImageURLs[index] {
			t.Errorf("ListAllImages() failed, misMatch: [%s, %s]\n", nameStr, TestImageURLs[index])
		}
	}
}

// 测试删除镜像
func TestRemoveImage(t *testing.T) {
	// 创建一个ImageManager
	im := &ImageManager{}

	// 遍历imageURLs，删除镜像
	for _, imageURL := range TestImageURLs {
		// 调用RemoveImage函数
		err := im.RemoveImage(imageURL)

		if err != nil {
			t.Errorf("RemoveImage() failed: %s\n", err)
		}
	}

	// 检查获取到的所有镜像数组是否长度为0
	images, err := im.ListAllImages()
	if err != nil {
		t.Errorf("ListAllImages() failed: %s\n", err)
	}
	if len(images) != 0 {
		t.Errorf("ListAllImages() failed: images is not empty")
	}
}

func TestPullImageWithPolicy(t *testing.T) {
	// 创建一个ImageManager
	im := &ImageManager{}

	// 遍历imageURLs
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImageWithPolicy(imageURL, ImagePullPolicyAlways)

		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}

	// 遍历imageURLs, 这次的策略是ImagePullPolicyIfNotPresent
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImageWithPolicy(imageURL, ImagePullPolicyIfNotPresent)

		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}

	// 遍历imageURLs,然后删除镜像
	for _, imageURL := range TestImageURLs {
		// 调用RemoveImage函数
		err := im.RemoveImage(imageURL)

		if err != nil {
			t.Errorf("RemoveImage() failed: %s\n", err)
		}
	}

	// 遍历imageURLs, 这次的策略是ImagePullPolicyIfNotPresent
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImageWithPolicy(imageURL, ImagePullPolicyIfNotPresent)

		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}

	// 遍历imageURLs, 这次的策略是ImagePullPolicyNever
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImageWithPolicy(imageURL, ImagePullPolicyNever)

		// 理论上这次应该是正常的
		if err != nil {
			t.Errorf("PullImage() failed: %s\n", err)
		}
		if image == nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}

	// 遍历imageURLs,然后删除镜像
	for _, imageURL := range TestImageURLs {
		// 调用RemoveImage函数
		err := im.RemoveImage(imageURL)

		if err != nil {
			t.Errorf("RemoveImage() failed: %s\n", err)
		}
	}

	// 遍历imageURLs, 这次的策略是ImagePullPolicyNever
	for _, imageURL := range TestImageURLs {
		// 调用PullImage函数
		image, err := im.PullImageWithPolicy(imageURL, ImagePullPolicyNever)

		// 理论上这次的error应该是"image not found"
		if err == nil || err.Error() != "image not found" {
			t.Errorf("PullImage() failed: %s\n", err)
		}

		if image != nil {
			t.Errorf("PullImage() failed: %s\n", image)
		}
	}

}
