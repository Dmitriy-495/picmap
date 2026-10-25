#!/usr/bin/env python3
import getpass
import bcrypt
import json

def hash_password(password: str) -> str:
    salt = bcrypt.gensalt()
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')

def main():
    print("=" * 50)
    print("🔐 Генератор хешей для PicMap")
    print("=" * 50)
    print("Пароли не будут отображаться при вводе.\n")

    users = []

    print("👤 Пользователь: tda")
    pwd1 = getpass.getpass("   Введи пароль: ")
    pwd1_confirm = getpass.getpass("   Повтори пароль: ")
    if pwd1 != pwd1_confirm:
        print("❌ Пароли не совпадают!")
        return
    users.append({
        "username": "tda",
        "email": "tda495@yandex.ru",
        "is_owner": True,
        "password_hash": hash_password(pwd1)
    })
    print("   ✅ Готово\n")

    print("👤 Пользователь: vra")
    pwd2 = getpass.getpass("   Введи пароль: ")
    pwd2_confirm = getpass.getpass("   Повтори пароль: ")
    if pwd2 != pwd2_confirm:
        print("❌ Пароли не совпадают!")
        return
    users.append({
        "username": "vra",
        "email": "r.vishlov@yandex.ru",
        "is_owner": False,
        "password_hash": hash_password(pwd2)
    })
    print("   ✅ Готово\n")

    config = {"superusers": users}
    print("=" * 50)
    print("📄 Готовый config.json:")
    print("=" * 50)
    print(json.dumps(config, indent=2, ensure_ascii=False))
    print("=" * 50)

    save = input("\n💾 Сохранить в backend/config.json? (y/n): ").strip().lower()
    if save == 'y':
        with open('backend/config.json', 'w', encoding='utf-8') as f:
            json.dump(config, f, indent=2, ensure_ascii=False)
        print("✅ Сохранено в backend/config.json")
    else:
        print("ℹ️  Скопируй JSON выше и вставь вручную.")

if __name__ == "__main__":
    main()
