import 'package:flutter/material.dart';
import '../services/subscription_service.dart';

class AddSubscriptionPage extends StatefulWidget {
  const AddSubscriptionPage({super.key});

  @override
  State<AddSubscriptionPage> createState() => _AddSubscriptionPageState();
}

class _AddSubscriptionPageState extends State<AddSubscriptionPage> {
  final _formKey = GlobalKey<FormState>();
  final _service = SubscriptionService();

  final _nameController = TextEditingController();
  final _categoryController = TextEditingController();
  final _amountController = TextEditingController();

  String _currency = 'THB';
  String _billingCycle = 'monthly';
  DateTime? _billingDate;
  bool _isTrial = false;

  Future<void> _pickBillingDate() async {
    final date = await showDatePicker(
      context: context,
      firstDate: DateTime.now(),
      lastDate: DateTime(2100),
      initialDate: DateTime.now(),
    );

    if (date != null) {
      setState(() => _billingDate = date);
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate() || _billingDate == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please fill all required fields')),
      );
      return;
    }

    try {
      await _service.createSubscription(
        name: _nameController.text,
        category: _categoryController.text,
        amount: double.parse(_amountController.text),
        currency: _currency,
        billingCycle: _billingCycle,
        billingDate: _billingDate!,
        status: 'active',
        isTrial: _isTrial,
      );

      Navigator.pop(context, true);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString())),
      );
    }
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Add Subscription')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: 'Name'),
                validator: (v) => v!.isEmpty ? 'Required' : null,
              ),

              TextFormField(
                controller: _categoryController,
                decoration: const InputDecoration(labelText: 'Category'),
              ),

              TextFormField(
                controller: _amountController,
                decoration: const InputDecoration(labelText: 'Amount'),
                keyboardType: TextInputType.number,
                validator: (v) =>
                    v!.isEmpty ? 'Enter amount' : null,
              ),

              const SizedBox(height: 16),

              DropdownButtonFormField(
                value: _currency,
                items: const [
                  DropdownMenuItem(value: 'THB', child: Text('THB')),
                  DropdownMenuItem(value: 'USD', child: Text('USD')),
                ],
                onChanged: (v) => setState(() => _currency = v!),
                decoration: const InputDecoration(labelText: 'Currency'),
              ),

              DropdownButtonFormField(
                value: _billingCycle,
                items: const [
                  DropdownMenuItem(value: 'monthly', child: Text('Monthly')),
                  DropdownMenuItem(value: 'yearly', child: Text('Yearly')),
                ],
                onChanged: (v) => setState(() => _billingCycle = v!),
                decoration:
                    const InputDecoration(labelText: 'Billing Cycle'),
              ),

              const SizedBox(height: 16),

              ListTile(
                title: Text(
                  _billingDate == null
                      ? 'Select Billing Date'
                      : _billingDate!.toLocal().toString().split(' ')[0],
                ),
                trailing: const Icon(Icons.calendar_today),
                onTap: _pickBillingDate,
              ),

              SwitchListTile(
                title: const Text('Trial'),
                value: _isTrial,
                onChanged: (v) => setState(() => _isTrial = v),
              ),
            ],
          ),
        ),
      ),

      bottomNavigationBar: Padding(
        padding: const EdgeInsets.all(16),
        child: ElevatedButton(
          onPressed: _submit,
          child: const Text('Save Subscription'),
        ),
      ),
    );
  }
}
