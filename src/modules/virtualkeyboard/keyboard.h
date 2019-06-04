 #include <linux/input-event-codes.h>

int  vkb_open();
int  vkb_set_info(int vendor, int product, const char *name);
int  vkb_add_key(int key);
int  vkb_register();
void vkb_close();

int vkb_emit_push(int key, int pressed);
int vkb_emit_sync();