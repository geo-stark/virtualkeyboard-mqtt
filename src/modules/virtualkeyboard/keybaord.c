#include <unistd.h>
#include <memory.h>
#include <errno.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <linux/uinput.h>
#include "keyboard.h"

static int fd_ = -1;

static int emit(int type, int code, int val)
{
	if (fd_ < 0)
		return EBADF;

	struct input_event ie;

	ie.type = type;
	ie.code = code;
	ie.value = val;
	/* timestamp values below are ignored */
	ie.time.tv_sec = 0;
	ie.time.tv_usec = 0;

	if (write(fd_, &ie, sizeof(ie)) < 0)
		return errno;
	return 0;
}

int vkb_open()
{
	int err;

	if (fd_ >= 0)
		return 0;

	fd_ = open("/dev/uinput", O_WRONLY | O_NONBLOCK);
	if (fd_ < 0)
		return errno;

	if (ioctl(fd_, UI_SET_EVBIT, EV_KEY) < 0)
		return errno;

	if ((err = vkb_set_info(0x1111, 0x2222, "virtual keyboard")) != 0)
		return err;

	return 0;
}

int vkb_set_info(int vendor, int product, const char *name)
{
	if (fd_ < 0)
		return EBADF;

	struct uinput_setup usetup;

	memset(&usetup, 0, sizeof(usetup));
	usetup.id.bustype = BUS_USB;
	usetup.id.vendor = vendor;
	usetup.id.product = product;
	strcpy(usetup.name, name);

	if (ioctl(fd_, UI_DEV_SETUP, &usetup) < 0)
		return errno;
	return 0;
}

int vkb_add_key(int key)
{
	if (fd_ < 0)
		return EBADF;

	if (ioctl(fd_, UI_SET_KEYBIT, key) < 0)
  		return errno;
	return 0;
}

int vkb_register()
{
	if (fd_ < 0)
		return EBADF;

	if (ioctl(fd_, UI_DEV_CREATE) < 0)
		return errno;

	/*
	* On UI_DEV_CREATE the kernel will create the device node for this
	* device. We are inserting a pause here so that userspace has time
	* to detect, initialize the new device, and can start listening to
	* the event, otherwise it will not notice the event we are about
	* to send. This pause is only needed in our example code!
	*/
	sleep(1);
	return 0;
}

void vkb_close()
{
  if (fd_ >= 0)
  {
	/*
	* Give userspace some time to read the events before we destroy the
	* device with UI_DEV_DESTOY.
	*/
	sleep(1);
	ioctl(fd_, UI_DEV_DESTROY);
	close(fd_);
	fd_ = -1;
  }
}

int vkb_emit_push(int key, int pressed)
{
	return emit(EV_KEY, key, pressed);
}

int vkb_emit_sync()
{
	return emit(EV_SYN, SYN_REPORT, 0);
}
